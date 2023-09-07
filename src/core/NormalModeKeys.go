package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func NormalModeEventsHandler(char rune, e *config.Editor) rune {
	if len(e.MotionBuffer) > 0 || utils.IsValidStartingChar(char, e.EditorMode) {
		e.MotionBuffer = append(e.MotionBuffer, char)
	}

	if len(e.MotionBuffer) > 4 {
		e.ClearMotionBuffer()
	}

	if len(e.MotionBuffer) > 1 {
		success := e.ExecuteMotion(string(e.MotionBuffer))
		if success {
			e.ClearMotionBuffer()
			return constants.INITIAL_REFRESH
		}
	} else {
		switch char {
		case ':':
			e.EditorMode = constants.EDITOR_MODE_FILE_BROWSER
			e.CurrentBuffer.StoredCx = e.Cx
			e.CurrentBuffer.StoredCy = e.Cy
			e.Cx = 0
			e.Cy = 0
			ReadHandler(e, e.RootDirectory)
			return constants.INITIAL_REFRESH
		case 'v':
			e.ClearMotionBuffer()
			e.SetMode(constants.EDITOR_MODE_VISUAL)
			e.HighlightSelection()
		case 'i':
			e.SetMode(constants.EDITOR_MODE_INSERT)
		case 'j':
			EditorMoveCursor(constants.ARROW_DOWN, e)
			return constants.ARROW_DOWN
		case 'k':
			EditorMoveCursor(constants.ARROW_UP, e)
			return constants.ARROW_UP
		case 'l':
			EditorMoveCursor(constants.ARROW_RIGHT, e)
			return constants.ARROW_RIGHT
		case 'h':
			EditorMoveCursor(constants.ARROW_LEFT, e)
			return constants.ARROW_LEFT
		case 'u':
			UndoAction(e)
		case utils.CTRL_KEY('r'):
			RedoAction(e)
		case constants.TAB_KEY:
			for i := 0; i < 4; i++ {
				EditorMoveCursor(constants.ARROW_RIGHT, e)
			}
		case constants.ENTER_KEY:
			EditorMoveCursor(constants.ARROW_DOWN, e)
			return constants.ARROW_DOWN
		case utils.CTRL_KEY(constants.QUIT_KEY):
			success := QuitKeyHandler(e)
			if !success {
				return char
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
		case '/':
			EditorFind(e)
			return constants.INITIAL_REFRESH
		case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
			DeleteHandler(e, char)
		case constants.PAGE_DOWN, constants.PAGE_UP:
			PageJumpHandler(e, char)
		}
		return char
	}
	return constants.INITIAL_REFRESH
}
