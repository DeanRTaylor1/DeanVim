package core

import (
	"github.com/deanrtaylor1/go-editor/config"
)

func RedoAction(e *config.Editor) {
	lastAction, success := e.CurrentBuffer.PopRedo()
	if !success {
		return
	}

	e.CurrentBuffer.AppendUndo(lastAction, e.UndoHistory)

	e.Cx = lastAction.Cx
	e.Cy = lastAction.Index
	e.CurrentBuffer.SliceIndex = lastAction.Cx - e.LineNumberWidth

	lastAction.RedoFunction()
}
