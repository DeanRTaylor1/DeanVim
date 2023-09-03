package core

import (
	"fmt"

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
		case 'n':
			e.CurrentBuffer.ClearSelection()
			e.SetMode(constants.EDITOR_MODE_NORMAL)
			return constants.INITIAL_REFRESH
		case 'j':
			// Adjust for line number offset
			adjustedCxEnd := e.CurrentBuffer.SelectedCxEnd - e.LineNumberWidth
			currentRowLength := len(e.GetCurrentRow().Chars)
			if adjustedCxEnd < currentRowLength {
				e.CurrentBuffer.SelectMoveRightBy(currentRowLength - adjustedCxEnd)
			}

			e.CurrentBuffer.SelectMoveDownBy(1)
			EditorMoveCursor(constants.ARROW_DOWN, e)
			return constants.ARROW_DOWN
		case 'k':
			// Retract the selection to the original cursor position
			diff := e.CurrentBuffer.SelectedCxEnd - (e.Cx - e.LineNumberWidth)
			if diff > 0 && e.Cy == e.CurrentBuffer.SelectedCyStart {
				config.LogToFile(fmt.Sprintf("MOVING BACK"))
				e.CurrentBuffer.SelectMoveLeftBy(diff - e.LineNumberWidth)
			}
			// Move the selection up by one line
			e.CurrentBuffer.SelectMoveUpBy(1)
			EditorMoveCursor(constants.ARROW_UP, e)
			return constants.ARROW_UP
		case 'l':
			if e.Cx-e.LineNumberWidth < len(e.GetCurrentRow().Chars) {
				e.CurrentBuffer.SelectMoveRightBy(1)
				EditorMoveCursor(constants.ARROW_RIGHT, e)
				return constants.ARROW_RIGHT
			}
			return constants.NO_OP
		case 'h':
			e.CurrentBuffer.SelectMoveLeftBy(1)
			EditorMoveCursor(constants.ARROW_LEFT, e)
			return constants.ARROW_LEFT
		case constants.TAB_KEY:
			for i := 0; i < 4; i++ {
				e.CurrentBuffer.SelectMoveRightBy(1)
				EditorMoveCursor(constants.ARROW_RIGHT, e)
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
