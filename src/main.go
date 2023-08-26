package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/term"

	"github.com/deanrtaylor1/go-editor/actions"
	"github.com/deanrtaylor1/go-editor/config"
	_ "github.com/deanrtaylor1/go-editor/highlighting"
)

/** file i/o **/

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

	err := enableRawMode(cfg)
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(syscall.Stdin), cfg.TerminalState)

	initEditor(cfg)

	if len(os.Args) >= 2 {
		err := actions.EditorOpen(cfg, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	actions.EditorSetStatusMessage(cfg, "HELP: CTRL-S = Save | Ctrl-Q = quit | Ctr-f = find")

	for {
		config.LogToFile(fmt.Sprintf("%d, %d, %d", cfg.Cx, cfg.SliceIndex, cfg.ColOff))
		actions.EditorRefreshScreen(cfg)
		actions.ProcessKeyPress(cfg.Reader, cfg)
	}
}
