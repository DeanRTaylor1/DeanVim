package main

import (
	"log"
	"os"
	"syscall"

	"golang.org/x/term"

	"github.com/deanrtaylor1/go-editor/actions"
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	_ "github.com/deanrtaylor1/go-editor/highlighting"
	"github.com/deanrtaylor1/go-editor/mappings"
)

func enableRawMode(cfg *config.EditorConfig) error {
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return err
	}
	cfg.TerminalState = oldState
	return nil
}

func initEditor(cfg *config.EditorConfig) {
	err := config.GetWindowSize(cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.ScreenRows -= 2
}

func main() {
	cfg := config.NewEditorConfig()
	motions := mappings.InitializeMotionMap(cfg)
	cfg.MotionMap = motions

	err := enableRawMode(cfg)
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(syscall.Stdin), cfg.TerminalState)

	initEditor(cfg)

	if len(os.Args) >= 2 {
		actions.ReadHandler(cfg, os.Args[1])
	}

	char := constants.INITIAL_REFRESH

	for {
		actions.EditorRefreshScreen(cfg, char)
		char = actions.ProcessKeyPress(cfg.Reader, cfg)

	}
}
