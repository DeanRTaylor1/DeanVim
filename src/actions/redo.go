package actions

import (
	"github.com/deanrtaylor1/go-editor/config"
)

func RedoAction(cfg *config.EditorConfig) {
	lastAction, success := cfg.CurrentBuffer.PopRedo()
	if !success {
		return
	}

	cfg.CurrentBuffer.AppendUndo(lastAction, cfg.UndoHistory)

	cfg.Cx = lastAction.Cx
	cfg.Cy = lastAction.Index
	cfg.CurrentBuffer.SliceIndex = lastAction.Cx - 5

	lastAction.RedoFunction()
}
