package mappings

import (
	"strings"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/core"
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
			OpenModal(e)
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

func OpenModal(e *config.Editor) {
	files, err := config.ListFiles(e.RootDirectory)
	if err != nil {
		config.LogToFile(err.Error())
		return
	}
	e.ModalOpen = !e.ModalOpen
	e.Modal.Data = files
	filesString := strings.Join(files, ", ")
	config.LogToFile("Files listed: " + filesString)
}
