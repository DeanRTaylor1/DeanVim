package core

import (
	"fmt"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func VisualModeKeyPressProcessor(char rune, cfg *config.Editor) rune {
	if len(cfg.MotionBuffer) > 0 || utils.IsValidStartingChar(char, cfg.EditorMode) {
		cfg.MotionBuffer = append(cfg.MotionBuffer, char)
	}

	if len(cfg.MotionBuffer) > 4 {
		cfg.ClearMotionBuffer()
	}

	config.LogToFile(fmt.Sprintf("MotionBuffer: %s", string(cfg.MotionBuffer)))

	if len(cfg.MotionBuffer) > 1 {
		success := cfg.ExecuteMotion(string(cfg.MotionBuffer))
		if success {
			cfg.ClearMotionBuffer()
			return constants.INITIAL_REFRESH
		}
	} else {
		switch char {
		case 'n':
			cfg.SetMode(constants.EDITOR_MODE_NORMAL)
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
		case constants.TAB_KEY:
			for i := 0; i < 4; i++ {
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
