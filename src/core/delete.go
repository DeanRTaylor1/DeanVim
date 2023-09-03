package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func EditorDelChar(cfg *config.Editor) {
	if cfg.Cy == cfg.CurrentBuffer.NumRows {
		return
	}
	if cfg.CurrentBuffer.SliceIndex == 0 && cfg.Cy == 0 {
		return
	}
	row := &cfg.CurrentBuffer.Rows[cfg.Cy]
	if cfg.CurrentBuffer.SliceIndex > 0 {
		if cfg.Cx-cfg.ColOff < cfg.LineNumberWidth {
			cfg.ColOff--
		}
		cfg.Cx--
		EditorRowDelChar(row, cfg.CurrentBuffer.SliceIndex-1, cfg)
		cfg.CurrentBuffer.SliceIndex--
	} else {
		cfg.Cx = cfg.CurrentBuffer.Rows[cfg.Cy-1].Length + cfg.LineNumberWidth
		cfg.CurrentBuffer.SliceIndex = cfg.CurrentBuffer.Rows[cfg.Cy-1].Length
		EditorDelRow(cfg)
		cfg.Cy--
	}
}

func EditorRowDelChar(row *config.Row, at int, cfg *config.Editor) {
	if at < 0 || at >= len(row.Chars) {
		return
	}
	if closingBracket, ok := constants.BracketPairs[rune(row.Chars[at])]; ok {
		// Check if the next character is the corresponding closing bracket
		if at+1 < len(row.Chars) && row.Chars[at+1] == byte(closingBracket) {
			// Delete the closing bracket along with the opening bracket
			copy(row.Chars[at:], row.Chars[at+2:])
			row.Chars = row.Chars[:len(row.Chars)-1]
		}
	}
	copy(row.Chars[at:], row.Chars[at+1:])
	row.Chars = row.Chars[:len(row.Chars)-1] // Access the Row field

	row.Length = len(row.Chars) // Update the length of the row
	EditorUpdateRow(row, cfg)
	cfg.CurrentBuffer.Dirty++
}

func EditorDelRow(cfg *config.Editor) {
	if cfg.Cy <= 0 || cfg.Cy >= cfg.CurrentBuffer.NumRows {
		return
	}

	mergeCurrentRowWithPrevious(cfg)
	updateRowIndicesFromCurrent(cfg)
	highlighting.ResetRowHighlights(-1, cfg)
	highlighting.SyntaxHighlightStateMachine(&cfg.CurrentBuffer.Rows[cfg.Cy-1], cfg)
	ResetRowTabs(cfg.Cy-1, cfg)
	cfg.CurrentBuffer.RemoveRowAtIndex(cfg.Cy)
	// deleteCurrentRow(cfg)
	cfg.CurrentBuffer.Dirty++
}

func ResetRowTabs(idx int, cfg *config.Editor) {
	row := &cfg.CurrentBuffer.Rows[idx]
	row.Tabs = make([]byte, row.Length)
}

func mergeCurrentRowWithPrevious(cfg *config.Editor) {
	prevRow := &cfg.CurrentBuffer.Rows[cfg.Cy-1]
	currentRow := &cfg.CurrentBuffer.Rows[cfg.Cy]

	prevRow.Chars = append(prevRow.Chars, currentRow.Chars...)
	prevRow.Length = len(prevRow.Chars)
}

func updateRowIndicesFromCurrent(cfg *config.Editor) {
	for i := cfg.Cy; i < len(cfg.CurrentBuffer.Rows); i++ {
		cfg.CurrentBuffer.Rows[i].Idx = i
	}
}