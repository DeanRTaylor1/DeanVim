package main

import (
	"log"
	"os"
	"syscall"

	"golang.org/x/term"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/core"
	_ "github.com/deanrtaylor1/go-editor/highlighting"
	"github.com/deanrtaylor1/go-editor/mappings"
)

func enableRawMode(cfg *config.Editor) error {
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return err
	}
	cfg.TerminalState = oldState
	return nil
}

func initEditor(cfg *config.Editor) {
	err := config.GetWindowSize(cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.ScreenRows -= 2
}

func main() {
	cfg := config.NewEditor()
	motions := mappings.InitializeMotionMap(cfg)
	cfg.MotionMap = motions

	err := enableRawMode(cfg)
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(syscall.Stdin), cfg.TerminalState)

	initEditor(cfg)

	if len(os.Args) >= 2 {
		core.ReadHandler(cfg, os.Args[1])
	}

	char := constants.INITIAL_REFRESH

	for {
		core.EditorRefreshScreen(cfg, char)
		char = core.ProcessKeyPress(cfg.Reader, cfg)

	}
}
