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

	"golang.org/x/term"
)

/** CONSTS **/
const VERSION = "0.0.1"

const (
	ARROW_LEFT  rune = 1000
	ARROW_RIGHT rune = 1001
	ARROW_UP    rune = 1002
	ARROW_DOWN  rune = 1003
	PAGE_UP     rune = 1004
	PAGE_DOWN   rune = 1005
	HOME_KEY    rune = 1006
	END_KEY     rune = 1007
	DEL_KEY     rune = 1008
	BACKSPACE   rune = 127
	QUIT_TIMES  int  = 3
	QUIT_KEY    rune = 'q'
	SAVE_KEY    rune = 's'
)

const TAB_STOP = 4

type Position struct {
	x int
	y int
}

func CTRL_KEY(ch rune) rune {
	return ch & 0x1f
}

/** DATA **/

func exitKey(key rune) bool {
	return key == 113
}

type erow struct {
	size  int
	chars []byte
}

type EditorConfig struct {
	cx            int
	cy            int
	screenRows    int
	screenCols    int
	terminalState *term.State
	numRows       int
	rows          [][]byte
	rowOff        int
	colOff        int
	fileName      string
	statusMsg     string
	statusMsgTime time.Time
	dirty         int
	quitTimes     int
}

func newErow() *erow {
	return &erow{
		size:  0,
		chars: []byte{},
	}
}

func newEditorConfig() *EditorConfig {
	return &EditorConfig{
		cx:            0,
		cy:            0,
		screenRows:    0,
		screenCols:    0,
		terminalState: nil,
		rows:          [][]byte{},
		numRows:       0,
		rowOff:        0,
		colOff:        0,
		fileName:      "[Not Selected]",
		statusMsg:     "",
		statusMsgTime: time.Time{},
		dirty:         0,
		quitTimes:     QUIT_TIMES,
	}
}

func getWindowSize(cfg *EditorConfig) error {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	cfg.screenCols = width
	cfg.screenRows = height
	return nil
}

/** file i/o **/

func replaceTabsWithSpaces(line []byte) []byte {
	var result []byte
	for _, b := range line {
		if b == '\t' {
			spacesNeeded := TAB_STOP - (len(result) % TAB_STOP)
			for j := 0; j < spacesNeeded; j++ {
				result = append(result, ' ')
			}
		} else {
			result = append(result, b)
		}
	}
	return result
}

func editorUpdateRow(line []byte, cfg *EditorConfig) {
	if cfg.cy < 1 {
		return
	}
	cfg.rows[cfg.cy] = line
}

func editorAppendRow(r []byte, cfg *EditorConfig) {
	convertedLine := replaceTabsWithSpaces(r)
	cfg.rows = append(cfg.rows, convertedLine)
	// editorUpdateRow(&convertedLine, cfg)
}

func editorDelRow(cfg *EditorConfig) {
	if cfg.cy <= 0 || cfg.cy >= cfg.numRows {
		return
	}

	// Append the current row to the previous one
	cfg.rows[cfg.cy-1] = append(cfg.rows[cfg.cy-1], cfg.rows[cfg.cy]...)

	// Delete the current row
	cfg.rows = append(cfg.rows[:cfg.cy], cfg.rows[cfg.cy+1:]...)

	cfg.numRows--
	cfg.dirty++
}

func editorRowInsertChar(row *[]byte, at int, char rune, cfg *EditorConfig) {
	if at < 0 || at > len(*row) {
		at = len(*row)
	}

	*row = append(*row, 0)
	copy((*row)[at+1:], (*row)[at:])
	(*row)[at] = byte(char)

	editorUpdateRow(*row, cfg)
	cfg.dirty++
}

func editorRowAppendString(cfg *EditorConfig) {
	if cfg.cy > 0 {
		cfg.rows[cfg.cy-1] = append(cfg.rows[cfg.cy-1], cfg.rows[cfg.cy]...)
	}
}

func editorRowDelChar(row *[]byte, at int, cfg *EditorConfig) {
	if at < 0 || at >= len(*row) {
		return
	}
	copy((*row)[at:], (*row)[at+1:])
	*row = (*row)[:len(*row)-1]

	editorUpdateRow(*row, cfg)
	cfg.dirty++
}

func editorInsertChar(char rune, cfg *EditorConfig) {
	if cfg.cy == cfg.numRows {
		editorAppendRow([]byte{}, cfg)
		cfg.numRows++
	}
	editorRowInsertChar(&cfg.rows[cfg.cy], cfg.cx, char, cfg)

	cfg.cx++
}

func editorDelChar(cfg *EditorConfig) {
	if cfg.cy == cfg.numRows {
		return
	}
	if cfg.cx == 0 && cfg.cy == 0 {
		return
	}

	row := &cfg.rows[cfg.cy]
	if cfg.cx > 0 {
		editorRowDelChar(row, cfg.cx-1, cfg)
		cfg.cx--
	} else {
		cfg.cx = len(cfg.rows[cfg.cy-1])
		editorDelRow(cfg)
		cfg.cy--
	}
}

func editorRowsToString(cfg *EditorConfig) string {
	var buffer strings.Builder
	for _, row := range cfg.rows {
		buffer.Write(row)
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func editorSave(cfg *EditorConfig) (string, error) {
	if cfg.fileName == "[Not Selected]" {
		return "", errors.New("no filename provided")
	}

	startTime := time.Now()
	content := editorRowsToString(cfg)

	file, err := os.OpenFile(cfg.fileName, os.O_RDWR|os.O_CREATE, 0644)
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
	numLines := len(cfg.rows)
	numBytes := len(content)
	message := fmt.Sprintf("\"%s\", %dL, %dB, %.3fms: written", cfg.fileName, numLines, numBytes, float64(elapsedTime.Nanoseconds())/1e6)

	cfg.dirty = 0

	return message, nil
}

func editorOpen(cfg *EditorConfig, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening file")
	}
	defer file.Close()
	cfg.fileName = file.Name()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		linelen := len(line)
		for linelen > 0 && (line[linelen-1] == '\n' || line[linelen-1] == '\r') {
			linelen--
		}
		row := []byte(line[:linelen])

		editorAppendRow(row, cfg)
		cfg.numRows++
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	cfg.dirty = 0

	return nil
}

func enableRawMode(cfg *EditorConfig) error {
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return err
	}
	cfg.terminalState = oldState
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
						return HOME_KEY, nil
					case '3':
						return DEL_KEY, nil
					case '4':
						return END_KEY, nil
					case '5':
						return PAGE_UP, nil
					case '6':
						return PAGE_DOWN, nil
					case '7':
						return HOME_KEY, nil
					case '8':
						return END_KEY, nil
					}
				}
			} else {
				switch seq[1] {
				case 'A':
					return ARROW_UP, nil // Up
				case 'B':
					return ARROW_DOWN, nil // Down
				case 'C':
					return ARROW_RIGHT, nil // Right
				case 'D':
					return ARROW_LEFT, nil // Left
				case 'H':
					return HOME_KEY, nil
				case 'F':
					return END_KEY, nil
				}
			}
		} else if seq[0] == 'O' {
			switch seq[1] {
			case 'H':
				return HOME_KEY, nil
			case 'F':
				return END_KEY, nil
			}
		}
		return '\x1b', nil
	} else {
		return char, nil // other keypresses
	}
}

func editorMoveCursor(key rune, cfg *EditorConfig) {
	row := []byte{}
	if cfg.cy < cfg.numRows {
		row = cfg.rows[cfg.cy]
	}
	// spacesNeeded := TAB_STOP - (cfg.cx % TAB_STOP)
	switch key {
	case rune(ARROW_LEFT):
		if cfg.cx != 0 {
			cfg.cx--
		} else if cfg.cy > 0 {
			cfg.cy--
			if cfg.cy < len(cfg.rows) {
				cfg.cx = len(cfg.rows[cfg.cy])
			}
		}
		break
	case rune(ARROW_RIGHT):
		if cfg.cx < len(cfg.rows[cfg.cy])-1 {
			cfg.cx++
		} else if cfg.cx == len(cfg.rows[cfg.cy]) && cfg.cy < len(cfg.rows)-1 {
			cfg.cy++
			cfg.cx = 0
		}
		break

	case rune(ARROW_DOWN):
		if cfg.cy < cfg.numRows {
			cfg.cy++
		}
		break
	case rune(ARROW_UP):
		if cfg.cy != 0 {
			cfg.cy--
		}
		break
	}

	// Reevaluate the row after processing the key input
	if cfg.cy < cfg.numRows {
		row = cfg.rows[cfg.cy]
	} else {
		row = []byte{}
	}

	// Check if the cursor is past the end of the current row
	rowLen := len(row)
	if cfg.cx > rowLen {
		cfg.cx = rowLen
	}
}

func processKeyPress(reader *bufio.Reader, cfg *EditorConfig) {
	char, err := readKey(reader)
	if err != nil {
		panic(err)
	}

	switch char {
	case '\r':
		break
	case CTRL_KEY(QUIT_KEY):
		if cfg.dirty > 0 && cfg.quitTimes > 0 {
			editorSetStatusMessage(cfg, "WARNING!!! File has unsaved changes. Press Ctrl-Q %d more times to quit.", cfg.quitTimes)
			cfg.quitTimes--
			return
		}
		fmt.Print("\033[2J")
		fmt.Print("\x1b[H")
		os.Exit(0)
		break
	case CTRL_KEY(SAVE_KEY):
		msg, err := editorSave(cfg)
		if err != nil {
			editorSetStatusMessage(cfg, "%s", err.Error())
			break
		}
		editorSetStatusMessage(cfg, "%s", msg)
		break
	case HOME_KEY:
		cfg.cx = 0
		break
	case END_KEY:
		cfg.cx = len(cfg.rows[cfg.cy])
		break
	case BACKSPACE, CTRL_KEY('h'), DEL_KEY:
		if char == DEL_KEY {
			editorMoveCursor(ARROW_RIGHT, cfg)
		}
		editorDelChar(cfg)
		break
	case PAGE_DOWN, PAGE_UP:

		times := cfg.screenRows
		for times > 0 {
			if char == PAGE_UP {
				editorMoveCursor(ARROW_UP, cfg)
			} else {
				editorMoveCursor(ARROW_DOWN, cfg)
			}
			times--
		}
	case CTRL_KEY('l'), '\x1b':
		break
	case rune(ARROW_DOWN), rune(ARROW_UP), rune(ARROW_RIGHT), rune(ARROW_LEFT):
		editorMoveCursor(char, cfg)
	default:
		editorInsertChar(char, cfg)
		break
	}
	cfg.quitTimes = QUIT_TIMES
}

func editorScroll(cfg *EditorConfig) {
	if cfg.cy < cfg.rowOff {
		cfg.rowOff = cfg.cy
	}
	if cfg.cy >= cfg.rowOff+cfg.screenRows {
		cfg.rowOff = cfg.cy - cfg.screenRows + 1
	}
	if cfg.cx < cfg.colOff {
		cfg.colOff = cfg.cx
	}
	if cfg.cx >= cfg.colOff+cfg.screenCols {
		cfg.colOff = cfg.cx - cfg.screenCols + 1
	}
}

func editorDrawRows(buffer *bytes.Buffer, cfg *EditorConfig) {
	screenRows := cfg.screenRows
	screenCols := cfg.screenCols
	for i := 0; i < screenRows; i++ {
		fileRow := i + cfg.rowOff
		if fileRow >= cfg.numRows {
			if cfg.numRows == 0 && i == screenRows/3 {
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
			rowLength := len(cfg.rows[fileRow]) - cfg.colOff
			if rowLength < 0 {
				rowLength = 0
			}
			if rowLength > screenCols {
				rowLength = screenCols
			}
			if cfg.colOff < len(cfg.rows[fileRow]) {
				buffer.Write(cfg.rows[fileRow][cfg.colOff : cfg.colOff+rowLength])
			} else {
				buffer.Write([]byte{})
			}
		}
		buffer.WriteString("\x1b[K")

		buffer.WriteString("\r\n")
	}
}

func editorDrawStatusBar(buf *bytes.Buffer, cfg *EditorConfig) {
	buf.WriteString("\x1b[7m")

	currentRow := cfg.cy + 1
	if currentRow > cfg.numRows {
		currentRow = cfg.numRows
	}

	dirty := ""
	if cfg.dirty > 0 {
		dirty = "(modified)"
	}

	status := fmt.Sprintf("%.20s - %d lines %s", cfg.fileName, cfg.numRows, dirty)
	rStatus := fmt.Sprintf("%d/%d, cx: %d, rowLen: %d", currentRow, cfg.numRows, cfg.cx, len(cfg.rows[cfg.cy]))

	rLen := len(rStatus)
	if len(status) > cfg.screenCols {
		status = status[:cfg.screenCols-rLen]
	}

	buf.WriteString(status)
	for i := len(status); i < cfg.screenCols-rLen; i++ {
		buf.WriteString(" ")
	}

	buf.WriteString(rStatus)
	buf.WriteString("\x1b[m")
	buf.WriteString("\r\n")
}

func editorDrawMessageBar(buf *bytes.Buffer, cfg *EditorConfig) {
	buf.WriteString("\x1b[K") // Clear the line
	msgLen := len(cfg.statusMsg)
	if msgLen > cfg.screenCols {
		msgLen = cfg.screenCols
	}
	if msgLen > 0 && time.Since(cfg.statusMsgTime).Seconds() < 5 {
		buf.WriteString(cfg.statusMsg[:msgLen]) // Write the message if within 5 seconds
	}
}

func editorRefreshScreen(cfg *EditorConfig) {
	editorScroll(cfg)
	var buffer bytes.Buffer

	buffer.WriteString("\x1b[?25l")
	buffer.WriteString("\x1b[H")

	editorDrawRows(&buffer, cfg)
	editorDrawStatusBar(&buffer, cfg)
	editorDrawMessageBar(&buffer, cfg)

	cursorPosition := fmt.Sprintf("\x1b[%d;%dH", (cfg.cy-cfg.rowOff)+1, (cfg.cx-cfg.colOff)+1)
	buffer.WriteString(cursorPosition)

	buffer.WriteString("\x1b[?25h")

	os.Stdout.Write(buffer.Bytes())
}

func editorSetStatusMessage(cfg *EditorConfig, format string, a ...interface{}) {
	cfg.statusMsg = fmt.Sprintf(format, a...)
	cfg.statusMsgTime = time.Now()
}

func initEditor(cfg *EditorConfig) {
	err := getWindowSize(cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.screenRows -= 2
}

func main() {
	cfg := newEditorConfig()

	err := enableRawMode(cfg)
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(syscall.Stdin), cfg.terminalState)

	initEditor(cfg)

	if len(os.Args) >= 2 {
		err := editorOpen(cfg, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	editorSetStatusMessage(cfg, "HELP: CTRL-S = Save | Ctrl-Q = quit")

	reader := bufio.NewReader(os.Stdin)

	for {
		editorRefreshScreen(cfg)
		processKeyPress(reader, cfg)
	}
}
