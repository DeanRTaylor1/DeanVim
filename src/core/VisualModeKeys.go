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
		case 'y':
			e.YankSelected()
			return constants.NO_OP
		case 'n':
			e.CurrentBuffer.ClearSelection()
			e.SetMode(constants.EDITOR_MODE_NORMAL)
			return constants.INITIAL_REFRESH
		case 'j':
			// Adjust for line number offset
			adjustedCxEnd := e.CurrentBuffer.SelectedCxEnd - e.LineNumberWidth
			currentRowLength := len(e.GetCurrentRow().Chars)
			if adjustedCxEnd < currentRowLength {
				e.SelectMoveRightBy(currentRowLength - adjustedCxEnd)
			}

			e.SelectMoveDownBy(1)
			EditorMoveCursor(constants.ARROW_DOWN, e)
			return constants.ARROW_DOWN
		case 'k':
			config.LogToFile(fmt.Sprintf("Initial e.Cy: %d, e.Cx: %d", e.Cy, e.Cx))
			if e.Cy == 0 {
				return constants.NO_OP
			}

			// Adjust for line number offset
			adjustedCxEnd := e.CurrentBuffer.SelectedCxEnd - e.LineNumberWidth
			adjustedCxStart := e.CurrentBuffer.SelectedCxStart - e.LineNumberWidth
			config.LogToFile(fmt.Sprintf("Adjusted CxEnd: %d, Adjusted CxStart: %d", adjustedCxEnd, adjustedCxStart))

			// Get the length of the row above the current row
			aboveRowLength := len(e.CurrentBuffer.Rows[e.Cy-1].Chars)
			config.LogToFile(fmt.Sprintf("Above Row Length: %d", aboveRowLength))

			if e.Cy == e.CurrentBuffer.SelectedCyEnd {
				config.LogToFile("Moving from the end")
				if e.Cy-1 == e.CurrentBuffer.SelectedCyStart {
					e.CurrentBuffer.SelectedCxEnd = e.Cx
					config.LogToFile("Updated SelectedCxEnd to cursor position")
				} else if adjustedCxEnd > aboveRowLength {
					e.SelectMoveLeftBy(adjustedCxEnd - aboveRowLength)
					config.LogToFile("Moved selection left")
				} else {
					e.CurrentBuffer.SelectedCxEnd = e.Cx
					config.LogToFile("Updated SelectedCxEnd to cursor position")
				}
			} else if e.Cy == e.CurrentBuffer.SelectedCyStart {
				config.LogToFile("Moving from the start")
				if adjustedCxStart > aboveRowLength {
					e.SelectMoveLeftBy(adjustedCxStart - aboveRowLength)
					config.LogToFile("Moved selection left")
				} else {
					e.CurrentBuffer.SelectedCxStart = e.Cx
					config.LogToFile("Updated SelectedCxStart to cursor position")
				}
			}

			// Move the selection up by one line
			e.SelectMoveUpBy(1)
			EditorMoveCursor(constants.ARROW_UP, e)
			config.LogToFile(fmt.Sprintf("Final e.Cy: %d, e.Cx: %d", e.Cy, e.Cx))
			return constants.ARROW_UP

		case 'l':
			if e.Cx-e.LineNumberWidth < len(e.GetCurrentRow().Chars) {
				e.SelectMoveRightBy(1)
				EditorMoveCursor(constants.ARROW_RIGHT, e)
				return constants.ARROW_RIGHT
			}
			return constants.NO_OP
		case 'h':

			if e.Cx > 5 {
				config.LogToFile(fmt.Sprintf("cfg.Cx: %d, SelectStart: %d, SelectEnd: %d", e.Cx, e.CurrentBuffer.SelectedCxStart, e.CurrentBuffer.SelectedCxEnd))
				e.SelectMoveLeftBy(1)
				EditorMoveCursor(constants.ARROW_LEFT, e)
				return constants.ARROW_LEFT
			}
			return constants.NO_OP

		case constants.TAB_KEY:
			for i := 0; i < 4; i++ {
				e.SelectMoveRightBy(1)
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
