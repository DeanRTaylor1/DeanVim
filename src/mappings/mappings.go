package mappings

import (
	"github.com/deanrtaylor1/go-editor/actions"
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func InitializeMotionMap(cfg *config.EditorConfig) map[string]func() {
	return map[string]func(){
		" pv": func() {
			cfg.EditorMode = constants.EDITOR_MODE_FILE_BROWSER
			cfg.CurrentBuffer.StoredCx = cfg.Cx
			cfg.CurrentBuffer.StoredCy = cfg.Cy
			cfg.Cx = 0
			cfg.Cy = 0
			actions.ReadHandler(cfg, cfg.RootDirectory)
		},
		// ... other mappings ...
	}
}
