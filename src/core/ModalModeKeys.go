package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func ModalModeEventsHandler(char rune, e *config.Editor) rune {
	switch char {
	case constants.ESCAPE_KEY:
		e.ModalOpen = false
	case constants.ARROW_DOWN:
		return constants.ARROW_DOWN
	case constants.ARROW_UP:
		return constants.ARROW_UP
	case constants.ARROW_RIGHT:
		return constants.ARROW_RIGHT
	case constants.ARROW_LEFT:
		EditorMoveCursor(constants.ARROW_LEFT, e)
		return constants.ARROW_LEFT
	case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
		DeleteHandler(e, char)
	default:
		insertCharModalInput(char, e)
	}
	return constants.INITIAL_REFRESH
}
