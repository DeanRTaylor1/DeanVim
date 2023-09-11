package main

import (
	"log"
	"os"
	"syscall"

	"golang.org/x/term"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/core"
	"github.com/deanrtaylor1/go-editor/mappings"
)

func enableRawMode(e *config.Editor) error {
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return err
	}
	e.TerminalState = oldState
	return nil
}

func initEditor(e *config.Editor) {
	err := config.GetWindowSize(e)
	if err != nil {
		log.Fatal(err)
	}
	e.ScreenRows -= 2
}

func main() {
	e := config.NewEditor()
	motions := mappings.InitializeMotionMap(e)
	e.MotionMap = motions

	err := enableRawMode(e)
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(syscall.Stdin), e.TerminalState)

	initEditor(e)

	if len(os.Args) >= 2 {
		core.ReadHandler(e, os.Args[1])
	}

	char := constants.INITIAL_REFRESH

	for {
		core.EditorRefreshScreen(e, char)
		char = core.EventHandlerMain(e.Reader, e)

	}
}
