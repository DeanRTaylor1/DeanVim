package core

import (
	"path/filepath"
	"strings"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/fuzzy"
	"github.com/deanrtaylor1/go-editor/grep"
	"github.com/deanrtaylor1/go-editor/utils"
)

func ModalModeEventsHandler(char rune, e *config.Editor) rune {
	switch char {
	case constants.ENTER_KEY:
		e.CacheCursorCoords()
		e.ResetCursorCoords()

		var fullPath string
		switch e.Modal.Type {
		case config.MODAL_TYPE_FUZZY:
			results, ok := e.Modal.Results.(fuzzy.Matches)
			if ok {
				fullPath = filepath.Join(e.RootDirectory, results[e.Modal.ItemIndex].Str)
			}
		default:
			results, ok := e.Modal.Results.(fuzzy.Matches)
			if ok {
				selectedResult := results[e.Modal.ItemIndex].Str
				parts := strings.Split(selectedResult, ":")
				if len(parts) >= 1 {
					filename := parts[0]
					fullPath = filepath.Join(e.RootDirectory, filename)
				}
			}
		}

		if fullPath == "" {
			// Handle error: could not determine the type of Results
			return constants.NO_OP
		}

		e.Modal.ItemIndex = 0
		e.Modal.DataRowOffset = 0
		e.Modal.SearchColOffset = 0
		e.Modal.CursorPosition = 0
		e.Modal.ModalInput = []byte{}
		e.ModalOpen = false
		e.Modal.ModalDrawn = false
		ReadHandler(e, fullPath)
		return constants.INITIAL_REFRESH

	case constants.ESCAPE_KEY:
		e.ModalOpen = false
		e.Modal.ModalDrawn = false
	case constants.ARROW_DOWN:
		switch e.Modal.Type {
		case config.MODAL_TYPE_FUZZY:
			data, ok := e.Modal.Data.(fuzzy.Matches)
			if ok && e.Modal.ItemIndex < len(data)-1 {
				e.Modal.ItemIndex++
				return constants.ARROW_DOWN
			}
		default:
			data, ok := e.Modal.Data.([]string)
			if ok && e.Modal.ItemIndex < len(data)-1 {
				e.Modal.ItemIndex++
				return constants.ARROW_DOWN
			}
		}
		return constants.NO_OP
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
		switch e.Modal.Type {
		case config.MODAL_TYPE_FUZZY:
			if e.Modal.CursorPosition > 0 {
				updateResults(e)
			} else {
				e.Modal.Results = e.Modal.Data
			}
		default:
			if e.Modal.CursorPosition > 0 {
				updateResults(e)
			} else {
				e.Modal.Results = []string{}
			}
		}

	default:
		insertCharModalInput(char, e)
		e.Modal.ResetToFirstItem()
		updateResults(e)
	}
	return constants.INITIAL_REFRESH
}

func updateResults(e *config.Editor) {
	switch e.Modal.Type {
	case config.MODAL_TYPE_FUZZY:
		matches := fuzzy.FindFrom(string(e.Modal.ModalInput), &e.Modal)
		e.Modal.Results = matches
	default:
		data, err := grep.RunGrep(string(e.Modal.ModalInput), e.RootDirectory)
		if err != nil {
			return
		}
		e.Modal.Data = data
		matches := fuzzy.FindFrom(string(e.Modal.ModalInput), &e.Modal)
		e.Modal.Results = matches
	}
}
