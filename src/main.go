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
	_ "github.com/deanrtaylor1/go-editor/highlighting"
)

var (
	lastMatch = -1
	direction = 1
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
				result = append(result, ' ')
			}
		} else {
			result = append(result, b)
		}
	}
	return result
}

func editorUpdateRow(line []byte, cfg *config.EditorConfig) {
	if cfg.Cy < 1 {
		return
	}
	cfg.Rows[cfg.Cy] = line
}

func editorInsertRow(r []byte, at int, cfg *config.EditorConfig) {
	if at < 0 || at > len(cfg.Rows) {
		convertedLine := replaceTabsWithSpaces(r)
		cfg.Rows = append(cfg.Rows, convertedLine)
		return
	}
	convertedLine := replaceTabsWithSpaces(r)
	cfg.Rows = append(cfg.Rows[:at], append([][]byte{convertedLine}, cfg.Rows[at:]...)...)

	// editorUpdateRow(&convertedLine, cfg)
}

func editorDelRow(cfg *config.EditorConfig) {
	if cfg.Cy <= 0 || cfg.Cy >= cfg.NumRows {
		return
	}

	// Append the current row to the previous one
	cfg.Rows[cfg.Cy-1] = append(cfg.Rows[cfg.Cy-1], cfg.Rows[cfg.Cy]...)

	// Delete the current row
	cfg.Rows = append(cfg.Rows[:cfg.Cy], cfg.Rows[cfg.Cy+1:]...)

	cfg.NumRows--
	cfg.Dirty++
}

func editorRowInsertChar(row *[]byte, at int, char rune, cfg *config.EditorConfig) {
	if at < 0 || at > len(*row) {
		at = len(*row)
	}

	*row = append(*row, 0)
	copy((*row)[at+1:], (*row)[at:])
	(*row)[at] = byte(char)

	editorUpdateRow(*row, cfg)
	cfg.Dirty++
}

// func editorRowAppendString(cfg *config.EditorConfig) {
// 	if cfg.Cy > 0 {
// 		cfg.Rows[cfg.Cy-1] = append(cfg.Rows[cfg.Cy-1], cfg.Rows[cfg.Cy]...)
// 	}
// }

func editorRowDelChar(row *[]byte, at int, cfg *config.EditorConfig) {
	if at < 0 || at >= len(*row) {
		return
	}
	copy((*row)[at:], (*row)[at+1:])
	*row = (*row)[:len(*row)-1]

	editorUpdateRow(*row, cfg)
	cfg.Dirty++
}

func editorInsertChar(char rune, cfg *config.EditorConfig) {
	if cfg.Cy == cfg.NumRows {
		editorInsertRow([]byte{}, -1, cfg)
		cfg.NumRows++
	}
	editorRowInsertChar(&cfg.Rows[cfg.Cy], cfg.Cx, char, cfg)

	cfg.Cx++
}

func editorInsertNewLine(cfg *config.EditorConfig) {
	if cfg.Cx == 0 {
		row := []byte{}
		at := cfg.Cy
		if cfg.Cy == 0 {
			at = cfg.Cy
		}
		editorInsertRow(row, at, cfg)
	} else {
		row := cfg.Rows[cfg.Cy]
		cfg.Rows[cfg.Cy] = row[:cfg.Cx]
		editorInsertRow(row[cfg.Cx:], cfg.Cy+1, cfg)
		cfg.Cx = 0
	}
	cfg.NumRows++
	cfg.Cy++
}

func editorDelChar(cfg *config.EditorConfig) {
	if cfg.Cy == cfg.NumRows {
		return
	}
	if cfg.Cx == 0 && cfg.Cy == 0 {
		return
	}

	row := &cfg.Rows[cfg.Cy]
	if cfg.Cx > 0 {
		editorRowDelChar(row, cfg.Cx-1, cfg)
		cfg.Cx--
	} else {
		cfg.Cx = len(cfg.Rows[cfg.Cy-1])
		editorDelRow(cfg)
		cfg.Cy--
	}
}

func editorRowsToString(cfg *config.EditorConfig) string {
	var buffer strings.Builder
	for _, row := range cfg.Rows {
		buffer.Write(row)
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
	numLines := len(cfg.Rows)
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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		linelen := len(line)
		for linelen > 0 && (line[linelen-1] == '\n' || line[linelen-1] == '\r') {
			linelen--
		}
		row := []byte(line[:linelen])

		editorInsertRow(row, -1, cfg)
		cfg.NumRows++
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	cfg.Dirty = 0

	return nil
}

func editorFindCallback(buf []rune, c rune, cfg *config.EditorConfig) {
	if c == '\r' || c == '\x1b' {
		lastMatch = -1
		direction = 1
		return
	} else if c == constants.ARROW_RIGHT || c == constants.ARROW_DOWN {
		direction = 1
	} else if c == constants.ARROW_LEFT || c == constants.ARROW_UP {
		direction = -1
	} else {
		lastMatch = -1
		direction = 1
	}

	if lastMatch == -1 {
		direction = 1
	}
	current := lastMatch
	for i := 0; i < cfg.NumRows; i++ {
		current += direction
		if current == -1 {
			current = cfg.NumRows - 1
		} else if current == cfg.NumRows {
			current = 0
		}

		row := cfg.Rows[current]
		matchIndex := strings.Index(string(row), string(buf))
		if matchIndex != -1 {
			lastMatch = current
			cfg.Cy = current
			cfg.Cx = matchIndex
			cfg.RowOff = cfg.NumRows
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
	if cfg.Cy < cfg.NumRows {
		row = cfg.Rows[cfg.Cy]
	}
	// spacesNeeded := TAB_STOP - (cfg.Cx % TAB_STOP)
	switch key {
	case rune(constants.ARROW_LEFT):
		if cfg.Cx != 0 {
			cfg.Cx--
		} else if cfg.Cy > 0 {
			cfg.Cy--
			if cfg.Cy < len(cfg.Rows) {
				cfg.Cx = len(cfg.Rows[cfg.Cy])
			}
		}
		break
	case rune(constants.SAVE_KEY):
		if cfg.Cy == cfg.NumRows {
			break
		}
		if cfg.Cx < len(cfg.Rows[cfg.Cy])-1 {
			cfg.Cx++
		} else if cfg.Cx == len(cfg.Rows[cfg.Cy]) && cfg.Cy < len(cfg.Rows)-1 {
			cfg.Cy++
			cfg.Cx = 0
		}
		break
	case rune(constants.ARROW_RIGHT):
		if cfg.Cy == cfg.NumRows {
			break
		}
		if cfg.Cx < len(cfg.Rows[cfg.Cy])-1 {
			cfg.Cx++
		} else if cfg.Cx == len(cfg.Rows[cfg.Cy]) && cfg.Cy < len(cfg.Rows)-1 {
			cfg.Cy++
			cfg.Cx = 0
		}
		break

	case rune(constants.ARROW_DOWN):
		if cfg.Cy < cfg.NumRows {
			cfg.Cy++
		}
		break
	case rune(constants.ARROW_UP):
		if cfg.Cy != 0 {
			cfg.Cy--
		}
		break
	}

	if cfg.Cy < cfg.NumRows {
		row = cfg.Rows[cfg.Cy]
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
		fmt.Print("\033[2J")
		fmt.Print("\x1b[H")
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
		if cfg.Cy == cfg.NumRows {
			break
		}
		cfg.Cx = len(cfg.Rows[cfg.Cy])
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
		if fileRow >= cfg.NumRows {
			if cfg.NumRows == 0 && i == screenRows/3 {
				welcome := "Go editor -- version 0.1"
				welcomelen := len(welcome)
				if welcomelen > screenCols {
					welcomelen = screenCols
				}
				padding := (screenCols - welcomelen) / 2
				if padding > 0 {
					buffer.WriteByte('~')
					padding--
				}
				for padding > 0 {
					buffer.WriteByte(' ')
					padding--
				}
				buffer.WriteString(welcome)
			} else {
				buffer.WriteByte('~')
			}
		} else {

			rowLength := len(cfg.Rows[fileRow]) - cfg.ColOff
			if rowLength < 0 {
				rowLength = 0
			}
			if rowLength > screenCols {
				rowLength = screenCols
			}
			if cfg.ColOff < len(cfg.Rows[fileRow]) {
				for j := 0; j < rowLength; j++ {
					c := cfg.Rows[fileRow][cfg.ColOff+j]
					if unicode.IsDigit(rune(c)) {
						buffer.WriteString("\x1b[31m")
						buffer.WriteByte(c)
						buffer.WriteString("\x1b[39m")
					} else {
						buffer.WriteByte(c)
					}
				}
			} else {
				buffer.Write([]byte{})
			}
		}
		buffer.WriteString("\x1b[K")

		buffer.WriteString("\r\n")
	}
}

func editorDrawStatusBar(buf *bytes.Buffer, cfg *config.EditorConfig) {
	buf.WriteString("\x1b[7m")

	currentRow := cfg.Cy + 1
	if currentRow > cfg.NumRows {
		currentRow = cfg.NumRows
	}

	dirty := ""
	if cfg.Dirty > 0 {
		dirty = "(modified)"
	}

	status := fmt.Sprintf("%.20s - %d lines %s", cfg.FileName, cfg.NumRows, dirty)
	rStatus := fmt.Sprintf("%d/%d", currentRow, cfg.NumRows)

	rLen := len(rStatus)
	if len(status) > cfg.ScreenCols {
		status = status[:cfg.ScreenCols-rLen]
	}

	buf.WriteString(status)
	for i := len(status); i < cfg.ScreenCols-rLen; i++ {
		buf.WriteString(" ")
	}

	buf.WriteString(rStatus)
	buf.WriteString("\x1b[m")
	buf.WriteString("\r\n")
}

func editorDrawMessageBar(buf *bytes.Buffer, cfg *config.EditorConfig) {
	buf.WriteString("\x1b[K") // Clear the line
	msgLen := len(cfg.StatusMsg)
	if msgLen > cfg.ScreenCols {
		msgLen = cfg.ScreenCols
	}
	if msgLen > 0 && time.Since(cfg.StatusMsgTime).Seconds() < 5 {
		buf.WriteString(cfg.StatusMsg[:msgLen]) // Write the message if within 5 seconds
	}
}

func editorRefreshScreen(cfg *config.EditorConfig) {
	editorScroll(cfg)
	var buffer bytes.Buffer

	buffer.WriteString("\x1b[?25l")
	buffer.WriteString("\x1b[H")

	editorDrawRows(&buffer, cfg)
	editorDrawStatusBar(&buffer, cfg)
	editorDrawMessageBar(&buffer, cfg)

	cursorPosition := fmt.Sprintf("\x1b[%d;%dH", (cfg.Cy-cfg.RowOff)+1, (cfg.Cx-cfg.ColOff)+1)
	buffer.WriteString(cursorPosition)

	buffer.WriteString("\x1b[?25h")

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
