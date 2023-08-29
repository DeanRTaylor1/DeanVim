package actions

import (
	"fmt"

	"github.com/deanrtaylor1/go-editor/config"
)

func RedoAction(cfg *config.EditorConfig) {
	lastAction, success := cfg.CurrentBuffer.PopRedo()
	config.LogToFile(fmt.Sprintf("%v %v", lastAction, success))
	if !success {
		return
	}

	cfg.CurrentBuffer.AppendUndo(lastAction, cfg.UndoHistory)

	cfg.Cx = lastAction.Cx
	cfg.Cy = lastAction.Index
	cfg.SliceIndex = lastAction.Cx - 5

	lastAction.RedoFunction()
}
