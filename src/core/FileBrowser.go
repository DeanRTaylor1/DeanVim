package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func FileBrowserModeKeyPressProcessor(char rune, cfg *config.Editor) rune {
	switch char {
	case 'R':
		if cfg.IsDir() {
			return constants.NO_OP
		}
		EditorRename(cfg)
		ReadHandler(cfg, cfg.CurrentDirectory)
		return constants.INITIAL_REFRESH
	case '%':
		EditorCreate(cfg)
		ReadHandler(cfg, cfg.CurrentDirectory)
		return constants.INITIAL_REFRESH
	case utils.CTRL_KEY('l'), constants.ESCAPE_KEY:
		cfg.ClearMotionBuffer()
		return constants.INITIAL_REFRESH
	case 'D':
		EditorDelete(cfg)
		ReadHandler(cfg, cfg.CurrentDirectory)
		return constants.INITIAL_REFRESH
	case 'j', constants.ARROW_DOWN:
		EditorMoveCursor(constants.ARROW_DOWN, cfg)
		return constants.ARROW_DOWN
	case 'k', constants.ARROW_UP:
		EditorMoveCursor(constants.ARROW_UP, cfg)
		return constants.ARROW_UP
	case 'l', constants.ARROW_LEFT:
		EditorMoveCursor(constants.ARROW_RIGHT, cfg)
		return constants.ARROW_RIGHT
	case 'h', constants.ARROW_RIGHT:
		EditorMoveCursor(constants.ARROW_LEFT, cfg)
		return constants.ARROW_LEFT
	case constants.ENTER_KEY:
		ReadHandler(cfg, cfg.FileBrowserItems[cfg.Cy-len(cfg.InstructionsLines())].Path)
		return constants.INITIAL_REFRESH
	case utils.CTRL_KEY(constants.QUIT_KEY):
		success := QuitKeyHandler(cfg)
		if !success {
			return char
		}
	case constants.HOME_KEY:
		HomeKeyHandler(cfg)
	case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
		DeleteHandler(cfg, char)
	default:
		return constants.NO_OP
	}
	return char
}
