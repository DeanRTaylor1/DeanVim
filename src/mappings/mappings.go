package mappings

import (
	"fmt"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/core"
	"github.com/deanrtaylor1/go-editor/fuzzy"
)

func InitializeMotionMap(e *config.Editor) map[string]func() {
	return map[string]func(){
		" pv": func() {
			GoToFileBrowser(e)
		},
		"yy": func() {
			yankLine(e)
		},
		" pf": func() {
			OpenFuzzyModal(e)
		},
		" ps": func() {
			OpenGrepModal(e)
		},
	}
}

func GoToFileBrowser(e *config.Editor) {
	e.EditorMode = constants.EDITOR_MODE_FILE_BROWSER
	e.CacheCursorCoords()
	e.ResetCursorCoords()
	core.ReadHandler(e, e.RootDirectory)
}

func yankLine(e *config.Editor) {
	e.EditorMode = constants.EDITOR_MODE_VISUAL
	e.HighlightLine()
	e.YankSelection()
	e.Yank.Type = config.LineWise
	e.ClearSelection()
	e.EditorMode = constants.EDITOR_MODE_NORMAL
}

func OpenFuzzyModal(e *config.Editor) {
	files, err := config.ListFiles(e.RootDirectory)
	if err != nil {
		config.LogToFile(err.Error())
		return
	}
	e.ModalOpen = !e.ModalOpen

	e.Modal = config.InitModal(config.MODAL_TYPE_FUZZY)

	data := make(fuzzy.Matches, len(files))
	for i, str := range files {
		data[i] = fuzzy.Match{
			Str:            str,
			MatchedIndexes: []int{}, // Empty because no characters are matched
		}
	}

	e.Modal.Data = data
	e.Modal.Results = data

	config.LogToFile(fmt.Sprintf("Initial Modal Results: %v", e.Modal.Results))
	config.LogToFile(fmt.Sprintf("Initial Modal Data: %v", e.Modal.Data))
}

func OpenGrepModal(e *config.Editor) {
	// e.ModalOpen = !e.ModalOpen
	//
	e.ModalOpen = !e.ModalOpen
	e.Modal = config.InitModal(config.MODAL_TYPE_GENERIC)
}
