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

	err := enableRawMode(cfg)
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(syscall.Stdin), cfg.TerminalState)

	initEditor(cfg)

	if len(os.Args) >= 2 {
		arg := os.Args[1]
		fileInfo, err := os.Stat(arg)
		if err != nil {
			log.Fatal(err)
		}

		if fileInfo.IsDir() {
			cfg.SetMode(constants.EDITOR_MODE_FILE_BROWSER)
			actions.DirectoryOpen(cfg, arg)
		} else if arg == "." {
			cfg.SetMode(constants.EDITOR_MODE_FILE_BROWSER)
			currentDir, err := os.Getwd()
			if err != nil {
				log.Fatal("Could not get current directory")
			}
			actions.DirectoryOpen(cfg, currentDir)
		} else {
			err := actions.EditorOpen(cfg, arg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	actions.EditorSetStatusMessage(cfg, "HELP: CTRL-S = Save | Ctrl-Q = quit | Ctr-f = find")

	char := constants.INITIAL_REFRESH

	for {
		actions.EditorRefreshScreen(cfg, char)
		char = actions.ProcessKeyPress(cfg.Reader, cfg)
	}
}
