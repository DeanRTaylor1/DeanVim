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
	}
}

func ChangeToNormalMode(e *config.Editor) {
	e.EditorMode = constants.EDITOR_MODE_FILE_BROWSER
	e.CacheCursorCoords()
	e.ResetCursorCoords()
	core.ReadHandler(e, e.RootDirectory)
}
