package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"syscall"
	"unicode"

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
)

const TAB_STOP = 4

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
			for i := 0; i < spacesNeeded; i++ {
				result = append(result, ' ')
			}
		} else {
			result = append(result, b)
		}
	}
	return result
}

func editorAppendRow(r []byte, cfg *EditorConfig) {
	convertedLine := replaceTabsWithSpaces(r)
	cfg.rows = append(cfg.rows, convertedLine)
}

func editorOpen(cfg *EditorConfig, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening file")
	}
	defer file.Close()

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
	switch key {
	case rune(ARROW_LEFT):
		if cfg.cx != 0 {
			cfg.cx--
		} else if cfg.cy > 0 {
			cfg.cy--
			cfg.cx = len(cfg.rows[cfg.cy])
		}
		break
	case rune(ARROW_RIGHT):
		if len(row) > 0 && cfg.cx < len(row) {
			cfg.cx++
		} else if cfg.cx == len(row) {
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
	case 17:
		fmt.Print("\033[2J")
		fmt.Print("\x1b[H")
		os.Exit(0)
		break
	case HOME_KEY:
		cfg.cx = 0
		break
	case END_KEY:
		cfg.cx = cfg.screenCols - 1
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
	case rune(ARROW_DOWN), rune(ARROW_UP), rune(ARROW_RIGHT), rune(ARROW_LEFT):
		editorMoveCursor(char, cfg)
	default:
		if unicode.IsControl(char) {
			fmt.Printf("%d\n", char)
		}
		fmt.Printf("%d ('%c')\n", char, char)
		break
	}
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
			buffer.Write(cfg.rows[fileRow][cfg.colOff : cfg.colOff+rowLength])
		}
		buffer.WriteString("\x1b[K")

		if i < screenRows-1 {
			buffer.WriteString("\r\n")
		}
	}
}

func editorRefreshScreen(cfg *EditorConfig) {
	editorScroll(cfg)
	var buffer bytes.Buffer

	buffer.WriteString("\x1b[?25l")
	buffer.WriteString("\x1b[H")

	editorDrawRows(&buffer, cfg)

	cursorPosition := fmt.Sprintf("\x1b[%d;%dH", (cfg.cy-cfg.rowOff)+1, (cfg.cx-cfg.colOff)+1)
	buffer.WriteString(cursorPosition)

	buffer.WriteString("\x1b[?25h")

	os.Stdout.Write(buffer.Bytes())
}

func initEditor(cfg *EditorConfig) {
	err := getWindowSize(cfg)
	if err != nil {
		log.Fatal(err)
	}
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

	reader := bufio.NewReader(os.Stdin)

	for {
		editorRefreshScreen(cfg)
		processKeyPress(reader, cfg)
	}
}
