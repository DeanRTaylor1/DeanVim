package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func FileBrowserEventsHandler(char rune, e *config.Editor) rune {
	switch char {
	case 'R':
		if e.IsDir() {
			return constants.NO_OP
		}
		EditorRename(e)
		ReadHandler(e, e.CurrentDirectory)
		return constants.INITIAL_REFRESH
	case '%':
		EditorCreate(e)
		ReadHandler(e, e.CurrentDirectory)
		return constants.INITIAL_REFRESH
	case utils.CTRL_KEY('l'), constants.ESCAPE_KEY:
		e.ClearMotionBuffer()
		return constants.INITIAL_REFRESH
	case 'D':
		EditorDelete(e)
		ReadHandler(e, e.CurrentDirectory)
		return constants.INITIAL_REFRESH
	case 'j', constants.ARROW_DOWN:
		EditorMoveCursor(constants.ARROW_DOWN, e)
		return constants.ARROW_DOWN
	case 'k', constants.ARROW_UP:
		EditorMoveCursor(constants.ARROW_UP, e)
		return constants.ARROW_UP
	case 'l', constants.ARROW_LEFT:
		EditorMoveCursor(constants.ARROW_RIGHT, e)
		return constants.ARROW_RIGHT
	case 'h', constants.ARROW_RIGHT:
		EditorMoveCursor(constants.ARROW_LEFT, e)
		return constants.ARROW_LEFT
	case constants.ENTER_KEY:
		ReadHandler(e, e.FileBrowserItems[e.Cy-len(e.InstructionsLines())].Path)
		return constants.INITIAL_REFRESH
	case utils.CTRL_KEY(constants.QUIT_KEY):
		success := QuitKeyHandler(e)
		if !success {
			return char
		}
	case constants.HOME_KEY:
		HomeKeyHandler(e)
	case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
		DeleteHandler(e, char)
	default:
		return constants.NO_OP
	}
	return char
}
