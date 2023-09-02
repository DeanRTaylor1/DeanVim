package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func UndoAction(cfg *config.Editor) {
	lastAction, success := cfg.CurrentBuffer.PopUndo()
	if !success {
		return
	}

	cfg.CurrentBuffer.AppendRedo(lastAction, cfg.UndoHistory)

	switch lastAction.ActionType {
	case constants.ACTION_UPDATE_ROW:
		cfg.CurrentBuffer.ReplaceRowAtIndex(lastAction.Index, lastAction.Row)
		cfg.Cy = lastAction.Index
		cfg.Cx = lastAction.Cx
		cfg.CurrentBuffer.SliceIndex = cfg.Cx - cfg.LineNumberWidth
	case constants.ACTION_APPEND_ROW_TO_PREVIOUS:
		prevRow, ok := lastAction.PrevRow.(config.Row)
		if !ok {
			return
		}
		cfg.CurrentBuffer.Rows[lastAction.Index-1] = prevRow
		cfg.CurrentBuffer.InsertRowAtIndex(lastAction.Index, lastAction.Row)
		cfg.Cx = lastAction.Cx
		cfg.Cy = lastAction.Index
		cfg.CurrentBuffer.SliceIndex = lastAction.Cx - cfg.LineNumberWidth
	case constants.ACTION_INSERT_ROW:
		cfg.CurrentBuffer.RemoveRowAtIndex(lastAction.Index)
		cfg.CurrentBuffer.ReplaceRowAtIndex(lastAction.Index, lastAction.Row)
		cfg.Cx = lastAction.Cx
		cfg.Cy = lastAction.Index
		cfg.CurrentBuffer.SliceIndex = cfg.Cx - cfg.LineNumberWidth
	case constants.ACTION_INSERT_CHAR_AT_EOF:
		cfg.CurrentBuffer.RemoveRowAtIndex(lastAction.Index)
		cfg.Cx = lastAction.Cx
		cfg.Cy = lastAction.Index
		cfg.CurrentBuffer.SliceIndex = 0
	}
}
