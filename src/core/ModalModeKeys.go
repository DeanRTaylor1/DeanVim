package core

import (
	"path/filepath"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/fuzzy"
	"github.com/deanrtaylor1/go-editor/utils"
)

func ModalModeEventsHandler(char rune, e *config.Editor) rune {
	switch char {
	case constants.ENTER_KEY:
		e.CacheCursorCoords()
		e.ResetCursorCoords()
		fullPath := filepath.Join(e.RootDirectory, e.Modal.Results[e.Modal.ItemIndex].Str)
		e.Modal.ItemIndex = 0
		e.Modal.DataRowOffset = 0
		e.Modal.SearchColOffset = 0
		e.Modal.CursorPosition = 0
		e.Modal.ModalInput = []byte{}
		e.ModalOpen = false
		ReadHandler(e, fullPath)
		return constants.INITIAL_REFRESH
	case constants.ESCAPE_KEY:
		e.ModalOpen = false
		e.Modal.Data = config.InitModal().Data
	case constants.ARROW_DOWN:
		if e.Modal.ItemIndex >= len(e.Modal.Data) {
			return constants.NO_OP
		}
		e.Modal.ItemIndex++
		return constants.ARROW_DOWN
	case constants.ARROW_UP:
		if e.Modal.ItemIndex <= 0 {
			return constants.NO_OP
		}
		e.Modal.ItemIndex--
		return constants.ARROW_UP
	case constants.ARROW_RIGHT:
		EditorMoveCursor(constants.ARROW_RIGHT, e)
		return constants.ARROW_RIGHT
	case constants.ARROW_LEFT:
		EditorMoveCursor(constants.ARROW_LEFT, e)
		return constants.ARROW_LEFT
	case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
		DeleteHandler(e, char)
		e.Modal.ResetToFirstItem()
		if e.Modal.CursorPosition > 0 {
			updateResults(e)
		} else {
			e.Modal.Results = e.Modal.Data
		}
	default:
		insertCharModalInput(char, e)
		e.Modal.ResetToFirstItem()
		updateResults(e)
	}
	return constants.INITIAL_REFRESH
}

func updateResults(e *config.Editor) {
	matches := fuzzy.FindFrom(string(e.Modal.ModalInput), &e.Modal)
	var filteredData []string
	for _, match := range matches {
		filteredData = append(filteredData, match.Str)
	}

	if len(filteredData) == 0 {
		return
	}

	// Update Modal.Data with the filtered data
	e.Modal.Results = matches
}
