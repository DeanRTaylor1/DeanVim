package core

import (
	"fmt"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func VisualModeEventsHandler(char rune, cfg *config.Editor) rune {
	if len(cfg.MotionBuffer) > 0 || utils.IsValidStartingChar(char, cfg.EditorMode) {
		cfg.MotionBuffer = append(cfg.MotionBuffer, char)
	}

	if len(cfg.MotionBuffer) > 4 {
		cfg.ClearMotionBuffer()
	}

	if len(cfg.MotionBuffer) > 1 {
		success := cfg.ExecuteMotion(string(cfg.MotionBuffer))
		if success {
			cfg.ClearMotionBuffer()
			return constants.INITIAL_REFRESH
		}
	} else {
		switch char {
		case 'n':
			cfg.CurrentBuffer.ClearSelection()
			cfg.SetMode(constants.EDITOR_MODE_NORMAL)
			return constants.INITIAL_REFRESH
		case 'j':
			// Adjust for line number offset
			adjustedCxEnd := cfg.CurrentBuffer.SelectedCxEnd - cfg.LineNumberWidth
			currentRowLength := len(cfg.GetCurrentRow().Chars)
			if adjustedCxEnd < currentRowLength {
				cfg.CurrentBuffer.SelectMoveRightBy(currentRowLength - adjustedCxEnd)
			}

			cfg.CurrentBuffer.SelectMoveDownBy(1)
			EditorMoveCursor(constants.ARROW_DOWN, cfg)
			return constants.ARROW_DOWN
		case 'k':
			// Retract the selection to the original cursor position
			diff := cfg.CurrentBuffer.SelectedCxEnd - (cfg.Cx - cfg.LineNumberWidth)
			if diff > 0 && cfg.Cy == cfg.CurrentBuffer.SelectedCyStart {
				config.LogToFile(fmt.Sprintf("MOVING BACK"))
				cfg.CurrentBuffer.SelectMoveLeftBy(diff - cfg.LineNumberWidth)
			}
			// Move the selection up by one line
			cfg.CurrentBuffer.SelectMoveUpBy(1)
			EditorMoveCursor(constants.ARROW_UP, cfg)
			return constants.ARROW_UP
		case 'l':
			if cfg.Cx-cfg.LineNumberWidth < len(cfg.GetCurrentRow().Chars) {
				cfg.CurrentBuffer.SelectMoveRightBy(1)
				EditorMoveCursor(constants.ARROW_RIGHT, cfg)
				return constants.ARROW_RIGHT
			}
			return constants.NO_OP
		case 'h':
			cfg.CurrentBuffer.SelectMoveLeftBy(1)
			EditorMoveCursor(constants.ARROW_LEFT, cfg)
			return constants.ARROW_LEFT
		case constants.TAB_KEY:
			for i := 0; i < 4; i++ {
				cfg.CurrentBuffer.SelectMoveRightBy(1)
				EditorMoveCursor(constants.ARROW_RIGHT, cfg)
			}
		case constants.ENTER_KEY:
			EditorMoveCursor(constants.ARROW_DOWN, cfg)
			return constants.ARROW_DOWN
		case constants.HOME_KEY:
			HomeKeyHandler(cfg)
		case constants.END_KEY:
			err := EndKeyHandler(cfg)
			if err != nil {
				config.LogToFile(err.Error())
				break
			}
		case constants.PAGE_DOWN, constants.PAGE_UP:
			PageJumpHandler(cfg, char)
		case constants.ESCAPE_KEY, utils.CTRL_KEY('l'):
			cfg.SetMode(constants.EDITOR_MODE_NORMAL)
		}
		return char
	}
	return constants.INITIAL_REFRESH
}
