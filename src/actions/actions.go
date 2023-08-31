package actions

import (
	"bufio"
	"fmt"
	"time"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func InsertModeKeyPressProcessor(char rune, cfg *config.EditorConfig) {
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
			return
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
		clearRedos = false
		break
	case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
		DeleteHandler(cfg, char)
	case constants.PAGE_DOWN, constants.PAGE_UP:
		PageJumpHandler(cfg, char)
		clearRedos = false
	case utils.CTRL_KEY('l'), constants.ESCAPE_KEY:
		cfg.SetMode(constants.EDITOR_MODE_NORMAL)
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
}

func NormalModeKeyPressProcessor(char rune, cfg *config.EditorConfig) rune {
	switch char {
	case 'i':
		cfg.SetMode(constants.EDITOR_MODE_INSERT)
	case 'j':
		EditorMoveCursor(constants.ARROW_DOWN, cfg)
		return constants.ARROW_DOWN
	case 'k':
		EditorMoveCursor(constants.ARROW_UP, cfg)
		return constants.ARROW_UP
	case 'l':
		EditorMoveCursor(constants.ARROW_RIGHT, cfg)
		return constants.ARROW_RIGHT
	case 'h':
		EditorMoveCursor(constants.ARROW_LEFT, cfg)
		return constants.ARROW_LEFT
	case 'u':
		UndoAction(cfg)
	case utils.CTRL_KEY('r'):
		RedoAction(cfg)
	case constants.TAB_KEY:
		for i := 0; i < 4; i++ {
			EditorMoveCursor(constants.ARROW_RIGHT, cfg)
		}
	case constants.ENTER_KEY:
		EditorMoveCursor(constants.ARROW_DOWN, cfg)
		return constants.ARROW_DOWN
	case utils.CTRL_KEY(constants.QUIT_KEY):
		success := QuitKeyHandler(cfg)
		if !success {
			return char
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
	}
	return char
}

func ProcessKeyPress(reader *bufio.Reader, cfg *config.EditorConfig) rune {
	char, err := ReadKey(reader)
	if err != nil {
		panic(err)
	}

	if cfg.EditorMode == constants.EDITOR_MODE_NORMAL || cfg.EditorMode == constants.EDITOR_MODE_FILE_BROWSER {
		char = NormalModeKeyPressProcessor(char, cfg)
	} else if cfg.EditorMode == constants.EDITOR_MODE_INSERT {
		InsertModeKeyPressProcessor(char, cfg)
	}
	return char
}

func EditorMoveCursor(key rune, cfg *config.EditorConfig) {
	if cfg.EditorMode == constants.EDITOR_MODE_FILE_BROWSER {
		switch key {
		case rune(constants.ARROW_LEFT):
			cfg.MoveCursorLeft()
		case rune(constants.ARROW_RIGHT):
			if cfg.Cy == len(cfg.FileBrowserItems)+5 {
				break
			}
			cfg.MoveCursorRight()
		case rune(constants.ARROW_DOWN):
			if cfg.Cy < len(cfg.FileBrowserItems)+5 {
				cfg.MoveCursorDown()
			}
		case rune(constants.ARROW_UP):
			if cfg.Cy >= 5 {
				cfg.MoveCursorUp()
			}
		}
		return
	}
	config.LogToFile(fmt.Sprintf("cfg.Cx: %d, cfg.SliceIndex: %d", cfg.Cx, cfg.SliceIndex))
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
