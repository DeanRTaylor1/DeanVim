package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func InsertModeKeyPressProcessor(char rune, cfg *config.Editor) {
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
		if IsClosingBracket(char) && cfg.GetCurrentRow().Length > cfg.CurrentBuffer.SliceIndex && IsClosingBracket(rune(cfg.GetCurrentRow().Chars[cfg.CurrentBuffer.SliceIndex])) {
			cfg.Cx++
			cfg.CurrentBuffer.SliceIndex++
		} else {
			InsertCharHandler(cfg, char)
		}
	}
	if clearRedos {
		cfg.ClearRedoStack()
	}
	cfg.QuitTimes = constants.QUIT_TIMES
}
