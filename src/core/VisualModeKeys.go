package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func VisualModeEventsHandler(char rune, e *config.Editor) rune {
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
		case 'y':
			e.YankSelected()
			return constants.NO_OP
		case 'n':
			e.SetMode(constants.EDITOR_MODE_NORMAL)
			e.ClearSelection()
			return constants.INITIAL_REFRESH
		case 'j':
			if e.Cy == len(e.CurrentBuffer.Rows)-1 {
				return constants.NO_OP
			}
			EditorMoveCursor(constants.ARROW_DOWN, e)
			e.MoveSelection()
			return constants.ARROW_DOWN

		case 'k':
			if e.Cy == 0 {
				return constants.NO_OP
			}

			EditorMoveCursor(constants.ARROW_UP, e)
			e.MoveSelection()
			return constants.ARROW_UP

		case 'l':
			if e.Cx-e.LineNumberWidth < len(e.GetCurrentRow().Chars) {
				EditorMoveCursor(constants.ARROW_RIGHT, e)
				e.MoveSelection()
				return constants.ARROW_RIGHT
			}
			return constants.NO_OP
		case 'h':
			if e.Cx > 5 {
				EditorMoveCursor(constants.ARROW_LEFT, e)
				e.MoveSelection()
				return constants.ARROW_LEFT
			}
			return constants.NO_OP
		case constants.TAB_KEY:
			for i := 0; i < 4; i++ {
				EditorMoveCursor(constants.ARROW_RIGHT, e)
				e.MoveSelection()
			}
		case constants.ENTER_KEY:
			EditorMoveCursor(constants.ARROW_DOWN, e)
			return constants.ARROW_DOWN
		case constants.HOME_KEY:
			HomeKeyHandler(e)
		case constants.END_KEY:
			err := EndKeyHandler(e)
			if err != nil {
				config.LogToFile(err.Error())
				break
			}
		case constants.PAGE_DOWN, constants.PAGE_UP:
			PageJumpHandler(e, char)
		case constants.ESCAPE_KEY, utils.CTRL_KEY('l'):
			e.SetMode(constants.EDITOR_MODE_NORMAL)
		}
		return char
	}
	return constants.INITIAL_REFRESH
}
