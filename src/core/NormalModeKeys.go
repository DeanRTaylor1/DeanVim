package core

import (
	"fmt"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func NormalModeKeyPressProcessor(char rune, cfg *config.Editor) rune {
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
		case ':':
			cfg.EditorMode = constants.EDITOR_MODE_FILE_BROWSER
			cfg.CurrentBuffer.StoredCx = cfg.Cx
			cfg.CurrentBuffer.StoredCy = cfg.Cy
			cfg.Cx = 0
			cfg.Cy = 0
			ReadHandler(cfg, cfg.RootDirectory)
			return constants.INITIAL_REFRESH
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
		case '/':
			EditorFind(cfg)
			return constants.INITIAL_REFRESH
		case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
			DeleteHandler(cfg, char)
		case constants.PAGE_DOWN, constants.PAGE_UP:
			PageJumpHandler(cfg, char)
		}
		return char
	}
	return constants.INITIAL_REFRESH
}
