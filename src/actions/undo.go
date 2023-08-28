package actions

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func UndoAction(cfg *config.EditorConfig) {
	lastAction, success := cfg.CurrentBuffer.PopUndo()
	if !success {
		return
	}

	cfg.CurrentBuffer.AppendRedo(lastAction, cfg.UndoHistory)

	switch lastAction.ActionType {
	case constants.ACTION_UPDATE_ROW:
		cfg.CurrentBuffer.Rows[lastAction.Index] = lastAction.Row
		cfg.Cy = lastAction.Index
		cfg.Cx = lastAction.Cx
		cfg.SliceIndex = lastAction.Row.Length
	case constants.ACTION_APPEND_ROW_TO_PREVIOUS:
		prevRow, ok := lastAction.PrevRow.(config.Row)
		if !ok {
			return
		}
		cfg.CurrentBuffer.Rows[lastAction.Index-1] = prevRow

		cfg.CurrentBuffer.Rows = append(cfg.CurrentBuffer.Rows[:lastAction.Index], append([]config.Row{lastAction.Row}, cfg.CurrentBuffer.Rows[lastAction.Index:]...)...)
		cfg.Cx = lastAction.Cx
		cfg.Cy = lastAction.Index
		cfg.SliceIndex = 0
	}
}
