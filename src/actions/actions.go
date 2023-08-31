package actions

import (
	"bufio"
	"fmt"
	"time"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func ProcessKeyPress(reader *bufio.Reader, cfg *config.EditorConfig) rune {
	char, err := ReadKey(reader)
	if err != nil {
		panic(err)
	}

	clearRedos := true

	switch char {
	case utils.CTRL_KEY('z'):
		UndoAction(cfg)
		clearRedos = false
	case utils.CTRL_KEY('y'):
		RedoAction(cfg)
		clearRedos = false
	case constants.TAB_KEY:
		TabKeyHandler(cfg)
	case constants.ENTER_KEY:
		EnterKeyHandler(cfg)
	case utils.CTRL_KEY(constants.QUIT_KEY):
		success := QuitKeyHandler(cfg)
		if !success {
			return constants.QUIT_KEY
		}
	case utils.CTRL_KEY(constants.SAVE_KEY):
		SaveKeyHandler(cfg)
	case constants.HOME_KEY:
		HomeKeyHandler(cfg)
	case constants.END_KEY:
		err := EndKeyHandler(cfg)
		if err != nil {
			config.LogToFile(err.Error())
			break
		}
	case utils.CTRL_KEY('f'):
		EditorFind(cfg)
		break
	case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
		DeleteHandler(cfg, char)
	case constants.PAGE_DOWN, constants.PAGE_UP:
		PageJumpHandler(cfg, char)
		clearRedos = false
	case utils.CTRL_KEY('l'), '\x1b':
		break
	case rune(constants.ARROW_DOWN), rune(constants.ARROW_UP), rune(constants.ARROW_RIGHT), rune(constants.ARROW_LEFT):
		EditorMoveCursor(char, cfg)
		clearRedos = false
	default:
		if IsClosingBracket(char) && cfg.GetCurrentRow().Length > cfg.SliceIndex && IsClosingBracket(rune(cfg.GetCurrentRow().Chars[cfg.SliceIndex])) {
			cfg.Cx++
			cfg.SliceIndex++
		} else {
			InsertCharHandler(cfg, char)
		}
	}
	if clearRedos {
		cfg.ClearRedoStack()
	}
	cfg.QuitTimes = constants.QUIT_TIMES
	return char
}

func EditorMoveCursor(key rune, cfg *config.EditorConfig) {
	row := []byte{}
	if cfg.Cy < cfg.CurrentBuffer.NumRows {
		row = cfg.CurrentBuffer.Rows[cfg.Cy].Chars
	}
	// spacesNeeded := TAB_STOP - (cfg.Cx % TAB_STOP)
	switch key {
	case rune(constants.ARROW_LEFT):
		if cfg.SliceIndex != 0 {
			cfg.MoveCursorLeft()
		} else if cfg.Cy > 0 && cfg.Cy < cfg.CurrentBuffer.NumRows {
			cfg.MoveCursorUp()
			cfg.Cx = (cfg.GetCurrentRow().Length) + cfg.LineNumberWidth
			cfg.SliceIndex = cfg.GetCurrentRow().Length
		}
	case rune(constants.ARROW_RIGHT):
		if cfg.Cy == cfg.CurrentBuffer.NumRows {
			break
		}
		if cfg.SliceIndex < (cfg.GetCurrentRow().Length) {
			cfg.MoveCursorRight()
		} else if cfg.Cx-cfg.LineNumberWidth >= cfg.GetCurrentRow().Length && cfg.Cy < len(cfg.CurrentBuffer.Rows)-1 {
			cfg.MoveCursorDown()
			cfg.Cx = cfg.LineNumberWidth
			cfg.SliceIndex = 0
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
	if cfg.SliceIndex > rowLen {
		cfg.Cx = rowLen + cfg.LineNumberWidth
		cfg.SliceIndex = rowLen
	}
}

func EditorSetStatusMessage(cfg *config.EditorConfig, format string, a ...interface{}) {
	cfg.StatusMsg = fmt.Sprintf(format, a...)
	cfg.StatusMsgTime = time.Now()
}
