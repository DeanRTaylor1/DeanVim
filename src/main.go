package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"syscall"
	"unicode"

	"golang.org/x/term"
)

/** DATA **/

func exitKey(key rune) bool {
	return key == 113
}

func editorDrawRows(input *string, screenRows int) {
	for y := 0; y < screenRows; y++ {
		*input += "~"
		if y < screenRows-1 {
			*input += "\r\n"
		}
	}
}

type EditorConfig struct {
	screenRows    int
	screenCols    int
	terminalState *term.State
}

func newEditorConfig() *EditorConfig {
	return &EditorConfig{
		screenRows:    0,
		screenCols:    0,
		terminalState: nil,
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
	return char, nil
}

func processKeyPress(reader *bufio.Reader) {
	char, err := readKey(reader)
	if err != nil {
		panic(err)
	}

	switch {
	case exitKey(char):
		fmt.Print("\033[2J")
		fmt.Print("\x1b[H")
		os.Exit(0)
		break
	case unicode.IsControl(char):
		fmt.Printf("%d\n", char)
		break
	default:
		fmt.Printf("%d ('%c')\n", char, char)
		break
	}
}

func editorRefreshScreen(cfg *EditorConfig) {
	var ab string
	ab += "\033[2J"
	ab += "\x1b[H"
	editorDrawRows(&ab, cfg.screenRows)
	ab += "\x1b[H"
	os.Stdout.WriteString(ab)
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

	reader := bufio.NewReader(os.Stdin)

	for {
		editorRefreshScreen(cfg)
		processKeyPress(reader)
	}
}
