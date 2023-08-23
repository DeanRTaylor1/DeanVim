package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
	"unicode"

	"golang.org/x/term"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
	_ "github.com/deanrtaylor1/go-editor/highlighting"
)

func CTRL_KEY(ch rune) rune {
	return ch & 0x1f
}

/** DATA **/

func getWindowSize(cfg *config.EditorConfig) error {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	cfg.ScreenCols = width
	cfg.ScreenRows = height
	return nil
}

/** file i/o **/

func replaceTabsWithSpaces(line []byte) []byte {
	var result []byte
	for _, b := range line {
		if b == '\t' {
			spacesNeeded := constants.TAB_STOP - (len(result) % constants.TAB_STOP)
			for j := 0; j < spacesNeeded; j++ {
				result = append(result, byte(constants.SPACE_RUNE))
			}
		} else {
			result = append(result, b)
		}
	}
	return result
}

func editorUpdateRow(row *config.Row, cfg *config.EditorConfig) {
	if cfg.Cy < 1 {
		return
	}
	cfg.CurrentBuffer.Rows[cfg.Cy].Chars = row.Chars
	cfg.CurrentBuffer.Rows[cfg.Cy].Length = row.Length
	highlighting.EditorUpdateSyntax(&cfg.CurrentBuffer.Rows[cfg.Cy], cfg)
}

func editorInsertRow(row *config.Row, at int, cfg *config.EditorConfig) {
	// Replace tabs with spaces
	convertedChars := replaceTabsWithSpaces(row.Chars)
	row.Chars = convertedChars
	row.Length = len(convertedChars)
	row.Idx = at // Set the index to the insertion point
	highlighting.EditorUpdateSyntax(row, cfg)

	if at < 0 || at >= len(cfg.CurrentBuffer.Rows) {
		// If at is outside the valid range, append the row to the end
		cfg.CurrentBuffer.Rows = append(cfg.CurrentBuffer.Rows, *row)
		return
	}

	// If at is within the valid range, insert the row at the specified position
	cfg.CurrentBuffer.Rows = append(cfg.CurrentBuffer.Rows[:at], append([]config.Row{*row}, cfg.CurrentBuffer.Rows[at:]...)...)

	// Update the Idx of the subsequent rows
	for i := at + 1; i < len(cfg.CurrentBuffer.Rows); i++ {
		config.LogToFile(fmt.Sprintf("I: %d", i))
		cfg.CurrentBuffer.Rows[i].Idx = i
	}
}

func editorDelRow(cfg *config.EditorConfig) {
	if cfg.Cy <= 0 || cfg.Cy >= cfg.CurrentBuffer.NumRows {
		return
	}

	// Append the current row's characters to the previous one
	cfg.CurrentBuffer.Rows[cfg.Cy-1].Chars = append(cfg.CurrentBuffer.Rows[cfg.Cy-1].Chars, cfg.CurrentBuffer.Rows[cfg.Cy].Chars...)
	cfg.CurrentBuffer.Rows[cfg.Cy-1].Length = len(cfg.CurrentBuffer.Rows[cfg.Cy-1].Chars) // Update the length of the previous row

	for i := cfg.Cy; i < len(cfg.CurrentBuffer.Rows); i++ {
		cfg.CurrentBuffer.Rows[i].Idx = i
	}

	highlighting.EditorUpdateSyntax(&cfg.CurrentBuffer.Rows[cfg.Cy-1], cfg)

	// Delete the current row
	cfg.CurrentBuffer.Rows = append(cfg.CurrentBuffer.Rows[:cfg.Cy], cfg.CurrentBuffer.Rows[cfg.Cy+1:]...)

	cfg.CurrentBuffer.NumRows-- // Update NumRows within CurrentBuffer
	cfg.Dirty++
}

func editorRowInsertChar(row *config.Row, at int, char rune, cfg *config.EditorConfig) {
	// ... existing code ...

	row.Chars = append(row.Chars, 0)
	copy(row.Chars[at+1:], row.Chars[at:])
	row.Chars[at] = byte(char)

	// Add similar logic for row.Highlighting
	row.Highlighting = append(row.Highlighting, constants.HL_NORMAL)
	copy(row.Highlighting[at+1:], row.Highlighting[at:])
	row.Highlighting[at] = constants.HL_NORMAL // or the appropriate value

	row.Length = len(row.Chars) // Update the length of the row

	// Call EditorUpdateSyntax here to ensure the highlighting is updated as well
	highlighting.EditorUpdateSyntax(row, cfg)

	cfg.Dirty++
}

func editorRowDelChar(row *config.Row, at int, cfg *config.EditorConfig) {
	if at < 0 || at >= len(row.Chars) {
		return
	}
	copy(row.Chars[at:], row.Chars[at+1:])
	row.Chars = row.Chars[:len(row.Chars)-1] // Access the Row field

	row.Length = len(row.Chars) // Update the length of the row
	editorUpdateRow(row, cfg)
	cfg.Dirty++
}

func editorInsertChar(char rune, cfg *config.EditorConfig) {
	if cfg.Cy == cfg.CurrentBuffer.NumRows {
		editorInsertRow(config.NewRow(), -1, cfg)
		cfg.CurrentBuffer.NumRows++
	}
	editorRowInsertChar(&cfg.CurrentBuffer.Rows[cfg.Cy], cfg.Cx, char, cfg)

	cfg.Cx++
}

func editorInsertNewLine(cfg *config.EditorConfig) {
	if cfg.Cx == 0 {
		newRow := config.NewRow()
		at := cfg.Cy
		editorInsertRow(newRow, at, cfg)
	} else {
		row := cfg.CurrentBuffer.Rows[cfg.Cy]
		cfg.CurrentBuffer.Rows[cfg.Cy].Chars = row.Chars[:cfg.Cx]
		cfg.CurrentBuffer.Rows[cfg.Cy].Length = len(cfg.CurrentBuffer.Rows[cfg.Cy].Chars)
		newRow := config.Row{Chars: row.Chars[cfg.Cx:]}
		editorInsertRow(&newRow, cfg.Cy+1, cfg)
		cfg.Cx = 0
	}

	cfg.CurrentBuffer.NumRows++ // Update NumRows within CurrentBuffer
	cfg.Cy++
}

func editorDelChar(cfg *config.EditorConfig) {
	if cfg.Cy == cfg.CurrentBuffer.NumRows {
		return
	}
	if cfg.Cx == 0 && cfg.Cy == 0 {
		return
	}

	row := &cfg.CurrentBuffer.Rows[cfg.Cy]
	if cfg.Cx > 0 {
		editorRowDelChar(row, cfg.Cx-1, cfg)
		cfg.Cx--
	} else {
		cfg.Cx = cfg.CurrentBuffer.Rows[cfg.Cy-1].Length // Access the Row field
		editorDelRow(cfg)
		cfg.Cy--
	}
}

func editorRowsToString(cfg *config.EditorConfig) string {
	var buffer strings.Builder
	for _, row := range cfg.CurrentBuffer.Rows {
		buffer.Write(row.Chars)
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func editorSave(cfg *config.EditorConfig) (string, error) {
	if cfg.FileName == "[Not Selected]" {
		return "", errors.New("no filename provided")
	}

	startTime := time.Now()
	content := editorRowsToString(cfg)

	file, err := os.OpenFile(cfg.FileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := file.Truncate(int64(len(content))); err != nil {
		return "", fmt.Errorf("failed to truncate file: %w", err)
	}

	n, err := file.WriteString(content)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}
	if n != len(content) {
		return "", errors.New("unexpected number of bytes written to file")
	}

	elapsedTime := time.Since(startTime) // End timing
	numLines := len(cfg.CurrentBuffer.Rows)
	numBytes := len(content)
	message := fmt.Sprintf("\"%s\", %dL, %dB, %.3fms: written", cfg.FileName, numLines, numBytes, float64(elapsedTime.Nanoseconds())/1e6)

	cfg.Dirty = 0

	return message, nil
}

func editorOpen(cfg *config.EditorConfig, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening file")
	}
	defer file.Close()
	cfg.FileName = file.Name()

	highlighting.EditorSelectSyntaxHighlight(cfg)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		linelen := len(line)
		for linelen > 0 && (line[linelen-1] == '\n' || line[linelen-1] == '\r') {
			linelen--
		}
		row := config.NewRow() // Create a new Row using the NewRow function
		row.Chars = []byte(line[:linelen])
		row.Length = linelen
		row.Idx = len(cfg.CurrentBuffer.Rows)
		config.LogToFile(fmt.Sprintf("Idx: %d", row.Idx))
		highlighting.EditorUpdateSyntax(row, cfg)

		editorInsertRow(row, -1, cfg)
		cfg.CurrentBuffer.NumRows++ // Update NumRows within CurrentBuffer
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	cfg.Dirty = 0

	return nil
}

func editorFindCallback(buf []rune, c rune, cfg *config.EditorConfig) {
	if len(cfg.CurrentBuffer.SearchState.SavedHl) > 0 {
		sl := cfg.CurrentBuffer.SearchState.SavedHlLine
		cfg.CurrentBuffer.Rows[sl].Highlighting = cfg.CurrentBuffer.SearchState.SavedHl
	}

	if c == '\r' || c == '\x1b' {
		cfg.CurrentBuffer.SearchState.LastMatch = -1
		cfg.CurrentBuffer.SearchState.Direction = 1
		return
	} else if c == constants.ARROW_RIGHT || c == constants.ARROW_DOWN {
		cfg.CurrentBuffer.SearchState.Direction = 1
	} else if c == constants.ARROW_LEFT || c == constants.ARROW_UP {
		cfg.CurrentBuffer.SearchState.Direction = -1
	} else {
		cfg.CurrentBuffer.SearchState.LastMatch = -1
		cfg.CurrentBuffer.SearchState.Direction = 1
	}

	if cfg.CurrentBuffer.SearchState.LastMatch == -1 {
		cfg.CurrentBuffer.SearchState.Direction = 1
	}
	current := cfg.CurrentBuffer.SearchState.LastMatch
	for i := 0; i < cfg.CurrentBuffer.NumRows; i++ {
		current += cfg.CurrentBuffer.SearchState.Direction
		if current == -1 {
			current = cfg.CurrentBuffer.NumRows - 1
		} else if current == cfg.CurrentBuffer.NumRows {
			current = 0
		}

		row := cfg.CurrentBuffer.Rows[current].Chars
		matchIndex := strings.Index(string(row), string(buf))
		if matchIndex != -1 {
			cfg.CurrentBuffer.SearchState.LastMatch = current
			cfg.Cy = current
			cfg.Cx = matchIndex
			cfg.RowOff = cfg.CurrentBuffer.NumRows

			cfg.CurrentBuffer.SearchState.SavedHlLine = current
			cfg.CurrentBuffer.SearchState.SavedHl = make([]byte, len(cfg.CurrentBuffer.Rows[cfg.Cy].Highlighting))
			copy(cfg.CurrentBuffer.SearchState.SavedHl, cfg.CurrentBuffer.Rows[cfg.Cy].Highlighting)

			for i := 0; i < len(buf); i++ {
				cfg.CurrentBuffer.Rows[cfg.Cy].Highlighting[matchIndex+i] = constants.HL_MATCH // Assuming HL_MATCH is defined
			}
			break
		}
	}
}

func editorFind(cfg *config.EditorConfig) {
	cx := cfg.Cx
	cy := cfg.Cy
	rowOff := cfg.RowOff
	colOff := cfg.ColOff

	query := editorPrompt("Search: (ESC to cancel)", editorFindCallback, cfg)

	if query == nil {
		cfg.Cx = cx
		cfg.Cy = cy
		cfg.RowOff = rowOff
		cfg.ColOff = colOff
	}
}

func enableRawMode(cfg *config.EditorConfig) error {
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return err
	}
	cfg.TerminalState = oldState
	return nil
}

func readKey(reader *bufio.Reader) (rune, error) {
	char, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	// readRune returns one byte so we check if that byte is an escape character
	// That means an arrow key could have been pressed so we replace that arrow with our keys for navigation
	// which copy vims mappings

	if char == '\x1b' {
		seq := make([]rune, 3)
		seq[0], _, err = reader.ReadRune()
		if err != nil {
			return '\x1b', nil
		}
		seq[1], _, err = reader.ReadRune()
		if err != nil {
			return '\x1b', nil
		}

		if seq[0] == '[' {
			if seq[1] >= '0' && seq[1] <= '9' {
				seq[2], _, err = reader.ReadRune()
				if err != nil {
					return '\x1b', nil
				}
				if seq[2] == '~' {
					switch seq[1] {
					case '1':
						return constants.HOME_KEY, nil
					case '3':
						return constants.DEL_KEY, nil
					case '4':
						return constants.END_KEY, nil
					case '5':
						return constants.PAGE_UP, nil
					case '6':
						return constants.PAGE_DOWN, nil
					case '7':
						return constants.HOME_KEY, nil
					case '8':
						return constants.END_KEY, nil
					}
				}
			} else {
				switch seq[1] {
				case 'A':
					return constants.ARROW_UP, nil // Up
				case 'B':
					return constants.ARROW_DOWN, nil // Down
				case 'C':
					return constants.ARROW_RIGHT, nil // Right
				case 'D':
					return constants.ARROW_LEFT, nil // Left
				case 'H':
					return constants.HOME_KEY, nil
				case 'F':
					return constants.END_KEY, nil
				}
			}
		} else if seq[0] == 'O' {
			switch seq[1] {
			case 'H':
				return constants.HOME_KEY, nil
			case 'F':
				return constants.END_KEY, nil
			}
		}
		return '\x1b', nil
	} else {
		return char, nil // other keypresses
	}
}

func editorPrompt(prompt string, cb func([]rune, rune, *config.EditorConfig), cfg *config.EditorConfig) []rune {
	buf := []rune{}
	for {
		editorSetStatusMessage(cfg, "%s", fmt.Sprintf("%s %s", prompt, string(buf)))
		editorRefreshScreen(cfg)
		c, err := readKey(cfg.Reader)
		if err != nil {
			log.Fatal(err)
		}

		if c == constants.DEL_KEY || c == CTRL_KEY('h') || c == constants.BACKSPACE {
			if len(buf) != 0 {
				buf = buf[:len(buf)-1]
				if cb != nil {
					cb(buf, c, cfg)
				}
			}
		} else if c == '\x1b' {
			editorSetStatusMessage(cfg, "")
			if cb != nil {
				cb(buf, c, cfg)
			}
			return nil
		} else if c == '\r' {
			if len(buf) != 0 {
				editorSetStatusMessage(cfg, "")
				if cb != nil {
					cb(buf, c, cfg)
				}
				return buf
			}
		} else if c != CTRL_KEY('c') && c < 128 {
			buf = append(buf, c)
		}

		if cb != nil {
			cb(buf, c, cfg)
		}
	}
}

func editorMoveCursor(key rune, cfg *config.EditorConfig) {
	row := []byte{}
	if cfg.Cy < cfg.CurrentBuffer.NumRows {
		row = cfg.CurrentBuffer.Rows[cfg.Cy].Chars
	}
	// spacesNeeded := TAB_STOP - (cfg.Cx % TAB_STOP)
	switch key {
	case rune(constants.ARROW_LEFT):
		if cfg.Cx != 0 {
			cfg.Cx--
		} else if cfg.Cy > 0 {
			cfg.Cy--
			if cfg.Cy < len(cfg.CurrentBuffer.Rows) {
				cfg.Cx = (cfg.CurrentBuffer.Rows[cfg.Cy].Length)
			}
		}
		break
	case rune(constants.SAVE_KEY):
		if cfg.Cy == cfg.CurrentBuffer.NumRows {
			break
		}
		if cfg.Cx < (cfg.CurrentBuffer.Rows[cfg.Cy].Length)-1 {
			cfg.Cx++
		} else if cfg.Cx == (cfg.CurrentBuffer.Rows[cfg.Cy].Length) && cfg.Cy < len(cfg.CurrentBuffer.Rows)-1 {
			cfg.Cy++
			cfg.Cx = 0
		}
		break
	case rune(constants.ARROW_RIGHT):
		if cfg.Cy == cfg.CurrentBuffer.NumRows {
			break
		}
		if cfg.Cx < (cfg.CurrentBuffer.Rows[cfg.Cy].Length) {
			cfg.Cx++
		} else if cfg.Cx == cfg.CurrentBuffer.Rows[cfg.Cy].Length && cfg.Cy < len(cfg.CurrentBuffer.Rows)-1 {
			cfg.Cy++
			cfg.Cx = 0
		}
		break

	case rune(constants.ARROW_DOWN):
		if cfg.Cy < cfg.CurrentBuffer.NumRows {
			cfg.Cy++
		}
		break
	case rune(constants.ARROW_UP):
		if cfg.Cy != 0 {
			cfg.Cy--
		}
		break
	}

	if cfg.Cy < cfg.CurrentBuffer.NumRows {
		row = cfg.CurrentBuffer.Rows[cfg.Cy].Chars
	} else {
		row = []byte{}
	}

	rowLen := len(row)
	if cfg.Cx > rowLen {
		cfg.Cx = rowLen
	}
}

func processKeyPress(reader *bufio.Reader, cfg *config.EditorConfig) {
	char, err := readKey(reader)
	if err != nil {
		panic(err)
	}

	switch char {
	case constants.ENTER_KEY:
		editorInsertNewLine(cfg)
		break
	case CTRL_KEY(constants.QUIT_KEY):
		if cfg.Dirty > 0 && cfg.QuitTimes > 0 {
			editorSetStatusMessage(cfg, "WARNING!!! File has unsaved changes. Press Ctrl-Q %d more times to quit.", cfg.QuitTimes)
			cfg.QuitTimes--
			return
		}
		fmt.Print(constants.ESCAPE_CLEAR_SCREEN)
		fmt.Print(constants.ESCAPE_MOVE_TO_HOME_POS)
		os.Exit(0)
		break
	case CTRL_KEY(constants.SAVE_KEY):
		msg, err := editorSave(cfg)
		if err != nil {
			editorSetStatusMessage(cfg, "%s", err.Error())
			break
		}
		editorSetStatusMessage(cfg, "%s", msg)
		break
	case constants.HOME_KEY:
		cfg.Cx = 0
		break
	case constants.END_KEY:
		if cfg.Cy == cfg.CurrentBuffer.NumRows {
			break
		}
		cfg.Cx = cfg.CurrentBuffer.Rows[cfg.Cy].Length
		break
	case CTRL_KEY('f'):
		editorFind(cfg)
		break
	case constants.BACKSPACE, CTRL_KEY('h'), constants.DEL_KEY:
		if char == constants.DEL_KEY {
			editorMoveCursor(constants.ARROW_RIGHT, cfg)
		}
		editorDelChar(cfg)
		break
	case constants.PAGE_DOWN, constants.PAGE_UP:
		times := cfg.ScreenRows
		for times > 0 {
			if char == constants.PAGE_UP {
				editorMoveCursor(constants.ARROW_UP, cfg)
			} else {
				editorMoveCursor(constants.ARROW_DOWN, cfg)
			}
			times--
		}
	case CTRL_KEY('l'), '\x1b':
		break
	case rune(constants.ARROW_DOWN), rune(constants.ARROW_UP), rune(constants.ARROW_RIGHT), rune(constants.ARROW_LEFT):
		editorMoveCursor(char, cfg)
	default:
		editorInsertChar(char, cfg)
		break
	}
	cfg.QuitTimes = constants.QUIT_TIMES
}

func editorScroll(cfg *config.EditorConfig) {
	if cfg.Cy < cfg.RowOff {
		cfg.RowOff = cfg.Cy
	}
	if cfg.Cy >= cfg.RowOff+cfg.ScreenRows {
		cfg.RowOff = cfg.Cy - cfg.ScreenRows + 1
	}
	if cfg.Cx < cfg.ColOff {
		cfg.ColOff = cfg.Cx
	}
	if cfg.Cx >= cfg.ColOff+cfg.ScreenCols {
		cfg.ColOff = cfg.Cx - cfg.ScreenCols + 1
	}
}

func editorDrawRows(buffer *bytes.Buffer, cfg *config.EditorConfig) {
	screenRows := cfg.ScreenRows
	screenCols := cfg.ScreenCols
	for i := 0; i < screenRows; i++ {
		fileRow := i + cfg.RowOff
		if fileRow >= cfg.CurrentBuffer.NumRows {
			if cfg.CurrentBuffer.NumRows == 0 && i == screenRows/3 {
				welcome := "Go editor -- version 0.1"
				welcomelen := len(welcome)
				if welcomelen > screenCols {
					welcomelen = screenCols
				}
				padding := (screenCols - welcomelen) / 2
				if padding > 0 {
					buffer.WriteByte(byte(constants.TILDE))
					padding--
				}
				for padding > 0 {
					buffer.WriteByte(byte(constants.SPACE_RUNE))
					padding--
				}
				buffer.WriteString(welcome)
			} else {
				buffer.WriteByte(byte(constants.TILDE))
			}
		} else {
			rowLength := cfg.CurrentBuffer.Rows[fileRow].Length - cfg.ColOff
			if rowLength < 0 {
				rowLength = 0
			}
			if rowLength > screenCols {
				rowLength = screenCols
			}
			if cfg.ColOff < cfg.CurrentBuffer.Rows[fileRow].Length {
				if len(cfg.CurrentBuffer.Rows[fileRow].Highlighting) < 1 {
					panic("HIGHLIGHTING EMPTY")
				}
				highlights := cfg.CurrentBuffer.Rows[fileRow].Highlighting
				cColor := -1
				for j := 0; j < rowLength; j++ {
					c := cfg.CurrentBuffer.Rows[fileRow].Chars[cfg.ColOff+j]
					hl := highlights[cfg.ColOff+j]

					if unicode.IsControl(rune(c)) {
						sym := '?'
						if c <= 26 {
							sym = rune(int(c) + int('@'))
						}
						buffer.WriteString("\x1b[7m")
						buffer.WriteRune(sym)
						buffer.WriteString("\x1b[m")
						if cColor != -1 {
							buffer.WriteString(fmt.Sprintf("\x1b[%dm", cColor))
						}
					} else if hl == constants.HL_NORMAL {
						if cColor != -1 {
							buffer.WriteString(constants.FOREGROUND_RESET)
							cColor = -1
						}
						buffer.WriteByte(c)
					} else {
						color := int(highlighting.EditorSyntaxToColor(hl))
						if color != cColor {
							buffer.WriteString(fmt.Sprintf("\x1b[%dm", color))
							cColor = color
						}
						buffer.WriteByte(c)
					}
				}
				buffer.WriteString(constants.FOREGROUND_RESET)
			} else {
				buffer.Write([]byte{})
			}

		}
		buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)

		buffer.WriteString(constants.ESCAPE_NEW_LINE)
	}
}

func editorDrawStatusBar(buf *bytes.Buffer, cfg *config.EditorConfig) {
	buf.WriteString("\x1b[7m")

	currentRow := cfg.Cy + 1
	if currentRow > cfg.CurrentBuffer.NumRows {
		currentRow = cfg.CurrentBuffer.NumRows
	}

	dirty := ""
	if cfg.Dirty > 0 {
		dirty = "(modified)"
	}

	status := fmt.Sprintf("%.20s - %d lines %s", cfg.FileName, cfg.CurrentBuffer.NumRows, dirty)
	rStatus := fmt.Sprintf("%s | %d/%d", cfg.CurrentBuffer.BufferSyntax.FileType, cfg.Cy+1, cfg.CurrentBuffer.NumRows)

	rLen := len(rStatus)
	if len(status) > cfg.ScreenCols {
		status = status[:cfg.ScreenCols-rLen]
	}

	buf.WriteString(status)
	for i := len(status); i < cfg.ScreenCols-rLen; i++ {
		buf.WriteString(" ")
	}

	buf.WriteString(rStatus)
	buf.WriteString(constants.ESCAPE_RESET_ATTRIBUTES)
	buf.WriteString(constants.ESCAPE_NEW_LINE)
}

func editorDrawMessageBar(buf *bytes.Buffer, cfg *config.EditorConfig) {
	buf.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END) // Clear the line
	msgLen := len(cfg.StatusMsg)
	if msgLen > cfg.ScreenCols {
		msgLen = cfg.ScreenCols
	}
	if msgLen > 0 && time.Since(cfg.StatusMsgTime).Seconds() < 5 {
		buf.WriteString(cfg.StatusMsg[:msgLen]) // Write the message if within 5 seconds
	}
}

func setCursorPos(x int, y int) string {
	return fmt.Sprintf(constants.ESCAPE_MOVE_TO_COORDS, x, y)
}

func editorRefreshScreen(cfg *config.EditorConfig) {
	editorScroll(cfg)
	var buffer bytes.Buffer

	buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
	buffer.WriteString(constants.ESCAPE_MOVE_TO_HOME_POS)

	editorDrawRows(&buffer, cfg)
	editorDrawStatusBar(&buffer, cfg)
	editorDrawMessageBar(&buffer, cfg)

	cursorPosition := setCursorPos((cfg.Cy-cfg.RowOff)+1, (cfg.Cx-cfg.ColOff)+1)
	buffer.WriteString(cursorPosition)

	buffer.WriteString(constants.ESCAPE_SHOW_CURSOR)

	os.Stdout.Write(buffer.Bytes())
}

func editorSetStatusMessage(cfg *config.EditorConfig, format string, a ...interface{}) {
	cfg.StatusMsg = fmt.Sprintf(format, a...)
	cfg.StatusMsgTime = time.Now()
}

func initEditor(cfg *config.EditorConfig) {
	err := getWindowSize(cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.ScreenRows -= 2
}

func main() {
	cfg := config.NewEditorConfig()

	err := enableRawMode(cfg)
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(syscall.Stdin), cfg.TerminalState)

	initEditor(cfg)

	if len(os.Args) >= 2 {
		err := editorOpen(cfg, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	editorSetStatusMessage(cfg, "HELP: CTRL-S = Save | Ctrl-Q = quit | Ctr-f = find")

	for {
		editorRefreshScreen(cfg)
		processKeyPress(cfg.Reader, cfg)
	}
}
