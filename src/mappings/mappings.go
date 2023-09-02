package mappings

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/core"
)

func InitializeMotionMap(cfg *config.Editor) map[string]func() {
	return map[string]func(){
		" pv": func() {
			ChangeToNormalMode(cfg)
		},
	}
}

func ChangeToNormalMode(cfg *config.Editor) {
	cfg.EditorMode = constants.EDITOR_MODE_FILE_BROWSER
	cfg.CacheCursorCoords()
	cfg.ResetCursorCoords()
	core.ReadHandler(cfg, cfg.RootDirectory)
}
