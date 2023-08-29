package actions

import (
	"fmt"

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
		config.LogToFile(fmt.Sprintf("%s, %s", cfg.GetCurrentRow().Chars, lastAction.Row.Chars))

		cfg.CurrentBuffer.ReplaceRowAtIndex(lastAction.Index, lastAction.Row)

		cfg.Cy = lastAction.Index
		cfg.Cx = lastAction.Cx
		cfg.SliceIndex = cfg.Cx - cfg.LineNumberWidth
	case constants.ACTION_APPEND_ROW_TO_PREVIOUS:
		prevRow, ok := lastAction.PrevRow.(config.Row)
		if !ok {
			return
		}
		cfg.CurrentBuffer.Rows[lastAction.Index-1] = prevRow
		cfg.CurrentBuffer.InsertRowAtIndex(lastAction.Index, lastAction.Row)
		cfg.Cx = lastAction.Cx
		cfg.Cy = lastAction.Index
		cfg.SliceIndex = lastAction.Cx - cfg.LineNumberWidth
	case constants.ACTION_INSERT_ROW:
		cfg.CurrentBuffer.RemoveRowAtIndex(lastAction.Index)
		cfg.CurrentBuffer.ReplaceRowAtIndex(lastAction.Index, lastAction.Row)
		cfg.Cx = lastAction.Cx
		cfg.Cy = lastAction.Index
		cfg.SliceIndex = cfg.Cx - cfg.LineNumberWidth
		config.LogToFile(fmt.Sprintf("%s, Cx: %d, Cy: %d, RLength: %d, ARLength: %d, slice index: %d", cfg.GetCurrentRow().Chars, cfg.Cx, cfg.Cy, cfg.GetCurrentRow().Length, len(cfg.GetCurrentRow().Chars), cfg.SliceIndex))
	}
}
