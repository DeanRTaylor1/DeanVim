package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func UndoAction(e *config.Editor) {
	lastAction, success := e.CurrentBuffer.PopUndo()
	if !success {
		return
	}

	e.CurrentBuffer.AppendRedo(lastAction, e.UndoHistory)

	switch lastAction.ActionType {
	case constants.ACTION_UPDATE_ROW:
		e.CurrentBuffer.ReplaceRowAtIndex(lastAction.Index, lastAction.Row)
		e.Cy = lastAction.Index
		e.Cx = lastAction.Cx
		e.CurrentBuffer.SliceIndex = e.Cx - e.LineNumberWidth
	case constants.ACTION_APPEND_ROW_TO_PREVIOUS:
		prevRow, ok := lastAction.PrevRow.(config.Row)
		if !ok {
			return
		}
		e.CurrentBuffer.Rows[lastAction.Index-1] = prevRow
		e.CurrentBuffer.InsertRowAtIndex(lastAction.Index, lastAction.Row)
		e.Cx = lastAction.Cx
		e.Cy = lastAction.Index
		e.CurrentBuffer.SliceIndex = lastAction.Cx - e.LineNumberWidth
	case constants.ACTION_INSERT_ROW:
		e.CurrentBuffer.RemoveRowAtIndex(lastAction.Index)
		e.CurrentBuffer.ReplaceRowAtIndex(lastAction.Index, lastAction.Row)
		e.Cx = lastAction.Cx
		e.Cy = lastAction.Index
		e.CurrentBuffer.SliceIndex = e.Cx - e.LineNumberWidth
	case constants.ACTION_INSERT_CHAR_AT_EOF:
		e.CurrentBuffer.RemoveRowAtIndex(lastAction.Index)
		e.Cx = lastAction.Cx
		e.Cy = lastAction.Index
		e.CurrentBuffer.SliceIndex = 0
	}
}
