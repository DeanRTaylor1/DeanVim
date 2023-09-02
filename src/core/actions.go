package core

import (
	"bufio"
	"fmt"
	"time"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func ProcessKeyPress(reader *bufio.Reader, cfg *config.Editor) rune {
	char, err := ReadKey(reader)
	if err != nil {
		panic(err)
	}

	if cfg.EditorMode == constants.EDITOR_MODE_NORMAL {
		char = NormalModeKeyPressProcessor(char, cfg)
	} else if cfg.EditorMode == constants.EDITOR_MODE_INSERT {
		InsertModeKeyPressProcessor(char, cfg)
	} else if cfg.IsBrowsingFiles() {
		char = FileBrowserModeKeyPressProcessor(char, cfg)
	}
	return char
}

func FileBrowserCursorMovements(key rune, cfg *config.Editor) {
	switch key {
	case rune(constants.ARROW_LEFT):
		if cfg.Cx <= 0 {
			return
		}
		cfg.MoveCursorLeft()
	case rune(constants.ARROW_RIGHT):
		if cfg.Cy == len(cfg.FileBrowserItems)+len(cfg.InstructionsLines()) {
			break
		}
		cfg.MoveCursorRight()
	case rune(constants.ARROW_DOWN):
		if cfg.Cy < len(cfg.FileBrowserItems)+len(cfg.InstructionsLines()) {
			cfg.MoveCursorDown()
		}
	case rune(constants.ARROW_UP):
		if cfg.Cy >= 5 {
			cfg.MoveCursorUp()
		}
	}
}

func EditorCursorMovements(key rune, cfg *config.Editor) {
	var row []byte = []byte{}
	if cfg.Cy < cfg.CurrentBuffer.NumRows {
		row = cfg.CurrentBuffer.Rows[cfg.Cy].Chars
	}
	// spacesNeeded := TAB_STOP - (cfg.Cx % TAB_STOP)
	switch key {
	case rune(constants.ARROW_LEFT):
		if cfg.CurrentBuffer.SliceIndex != 0 {
			cfg.MoveCursorLeft()
		} else if cfg.Cy > 0 && cfg.Cy < cfg.CurrentBuffer.NumRows {
			cfg.MoveCursorUp()
			cfg.Cx = (cfg.GetCurrentRow().Length) + cfg.LineNumberWidth
			cfg.CurrentBuffer.SliceIndex = cfg.GetCurrentRow().Length
		}
	case rune(constants.ARROW_RIGHT):
		if cfg.Cy == cfg.CurrentBuffer.NumRows {
			break
		}
		if cfg.CurrentBuffer.SliceIndex < (cfg.GetCurrentRow().Length) {
			cfg.MoveCursorRight()
		} else if cfg.Cx-cfg.LineNumberWidth >= cfg.GetCurrentRow().Length && cfg.Cy < len(cfg.CurrentBuffer.Rows)-1 {
			cfg.MoveCursorDown()
			cfg.Cx = cfg.LineNumberWidth
			cfg.CurrentBuffer.SliceIndex = 0
		}
	case rune(constants.ARROW_DOWN):
		if cfg.Cy < cfg.CurrentBuffer.NumRows {
			cfg.MoveCursorDown()
		}
	case rune(constants.ARROW_UP):
		if cfg.Cy != 0 {
			cfg.MoveCursorUp()
		}
	}

	if cfg.Cy < cfg.CurrentBuffer.NumRows {
		row = cfg.CurrentBuffer.Rows[cfg.Cy].Chars
	} else {
		row = []byte{}
	}

	rowLen := len(row)
	if cfg.CurrentBuffer.SliceIndex > rowLen {
		cfg.Cx = rowLen + cfg.LineNumberWidth
		cfg.CurrentBuffer.SliceIndex = rowLen
	}
}

func EditorMoveCursor(key rune, cfg *config.Editor) {
	if cfg.IsBrowsingFiles() {
		FileBrowserCursorMovements(key, cfg)
	} else {
		EditorCursorMovements(key, cfg)
	}
}

func EditorSetStatusMessage(cfg *config.Editor, format string, a ...interface{}) {
	cfg.StatusMsg = fmt.Sprintf(format, a...)
	cfg.StatusMsgTime = time.Now()
}
