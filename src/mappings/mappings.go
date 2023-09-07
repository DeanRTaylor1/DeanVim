package mappings

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/core"
)

func InitializeMotionMap(e *config.Editor) map[string]func() {
	return map[string]func(){
		" pv": func() {
			ChangeToNormalMode(e)
		},
		"yy": func() {
			e.EditorMode = constants.EDITOR_MODE_VISUAL
			e.HighlightLine()
			e.YankSelection()
			e.Yank.Type = config.LineWise
			e.ClearSelection()
			e.EditorMode = constants.EDITOR_MODE_NORMAL
		},
	}
}

func ChangeToNormalMode(e *config.Editor) {
	e.EditorMode = constants.EDITOR_MODE_FILE_BROWSER
	e.CacheCursorCoords()
	e.ResetCursorCoords()
	core.ReadHandler(e, e.RootDirectory)
}
