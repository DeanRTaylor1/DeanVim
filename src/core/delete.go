package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func EditorDelChar(e *config.Editor) {
	if e.Cy == e.CurrentBuffer.NumRows {
		return
	}
	if e.CurrentBuffer.SliceIndex == 0 && e.Cy == 0 {
		return
	}
	row := &e.CurrentBuffer.Rows[e.Cy]
	if e.CurrentBuffer.SliceIndex > 0 {
		if e.Cx-e.ColOff < e.LineNumberWidth {
			e.ColOff--
		}
		e.Cx--
		EditorRowDelChar(row, e.CurrentBuffer.SliceIndex-1, e)
		e.CurrentBuffer.SliceIndex--
	} else {
		e.Cx = e.CurrentBuffer.Rows[e.Cy-1].Length + e.LineNumberWidth
		e.CurrentBuffer.SliceIndex = e.CurrentBuffer.Rows[e.Cy-1].Length
		EditorDelRow(e)
		e.Cy--
	}
}

func EditorRowDelChar(row *config.Row, at int, e *config.Editor) {
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
	EditorUpdateRow(row, e)
	e.CurrentBuffer.Dirty++
}

func EditorDelRow(e *config.Editor) {
	if e.Cy <= 0 || e.Cy >= e.CurrentBuffer.NumRows {
		return
	}

	mergeCurrentRowWithPrevious(e)
	updateRowIndicesFromCurrent(e)
	highlighting.ResetRowHighlights(-1, e)
	highlighting.SyntaxHighlightStateMachine(&e.CurrentBuffer.Rows[e.Cy-1], e)
	ResetRowTabs(e.Cy-1, e)
	e.CurrentBuffer.RemoveRowAtIndex(e.Cy)
	// deleteCurrentRow(e)
	e.CurrentBuffer.Dirty++
}

func ResetRowTabs(idx int, e *config.Editor) {
	row := &e.CurrentBuffer.Rows[idx]
	row.Tabs = make([]byte, row.Length)
}

func mergeCurrentRowWithPrevious(e *config.Editor) {
	prevRow := &e.CurrentBuffer.Rows[e.Cy-1]
	currentRow := &e.CurrentBuffer.Rows[e.Cy]

	prevRow.Chars = append(prevRow.Chars, currentRow.Chars...)
	prevRow.Length = len(prevRow.Chars)
}

func updateRowIndicesFromCurrent(e *config.Editor) {
	for i := e.Cy; i < len(e.CurrentBuffer.Rows); i++ {
		e.CurrentBuffer.Rows[i].Idx = i
	}
}
