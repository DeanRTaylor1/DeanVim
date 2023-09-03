package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func InsertModeEventsHandler(char rune, e *config.Editor) {
	clearRedos := true

	switch char {
	case utils.CTRL_KEY('z'):
		UndoAction(e)
		clearRedos = false
	case utils.CTRL_KEY('y'):
		RedoAction(e)
		clearRedos = false
	case constants.TAB_KEY:
		TabKeyHandler(e)
	case constants.ENTER_KEY:
		EnterKeyHandler(e)
	case utils.CTRL_KEY(constants.QUIT_KEY):
		success := QuitKeyHandler(e)
		if !success {
			return
		}
	case utils.CTRL_KEY(constants.SAVE_KEY):
		SaveKeyHandler(e)
	case constants.HOME_KEY:
		HomeKeyHandler(e)
	case constants.END_KEY:
		err := EndKeyHandler(e)
		if err != nil {
			config.LogToFile(err.Error())
			break
		}
	case utils.CTRL_KEY('f'):
		EditorFind(e)
		clearRedos = false
	case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
		DeleteHandler(e, char)
	case constants.PAGE_DOWN, constants.PAGE_UP:
		PageJumpHandler(e, char)
		clearRedos = false
	case utils.CTRL_KEY('l'), constants.ESCAPE_KEY:
		e.SetMode(constants.EDITOR_MODE_NORMAL)
	case rune(constants.ARROW_DOWN), rune(constants.ARROW_UP), rune(constants.ARROW_RIGHT), rune(constants.ARROW_LEFT):
		EditorMoveCursor(char, e)
		clearRedos = false
	default:
		if IsClosingBracket(char) && e.GetCurrentRow().Length > e.CurrentBuffer.SliceIndex && IsClosingBracket(rune(e.GetCurrentRow().Chars[e.CurrentBuffer.SliceIndex])) {
			e.Cx++
			e.CurrentBuffer.SliceIndex++
		} else {
			InsertCharHandler(e, char)
		}
	}
	if clearRedos {
		e.ClearRedoStack()
	}
	e.QuitTimes = constants.QUIT_TIMES
}
