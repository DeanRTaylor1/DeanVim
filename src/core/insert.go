package core

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func editorRowInsertChar(row *config.Row, at int, char rune, e *config.Editor) {
	row.Chars = append(row.Chars, 0)
	copy(row.Chars[at+1:], row.Chars[at:])
	row.Chars[at] = byte(char)

	row.Length = len(row.Chars)
	row.Highlighting = make([]byte, row.Length)
	highlighting.Fill(row.Highlighting, constants.HL_NORMAL)
	row.Tabs = make([]byte, row.Length)

	highlighting.SyntaxHighlightStateMachine(row, e)

	e.CurrentBuffer.Dirty++
}

func EditorInsertChar(char rune, e *config.Editor) {
	if e.Cy == e.CurrentBuffer.NumRows {
		EditorInsertRow(config.NewRow(), -1, e)
		e.CurrentBuffer.NumRows++
	}
	editorRowInsertChar(&e.CurrentBuffer.Rows[e.Cy], e.CurrentBuffer.SliceIndex, char, e)

	e.MoveCursorRight()
}

func EditorInsertNewLine(e *config.Editor) {
	row := e.CurrentBuffer.Rows[e.Cy]
	isBetweenBrackets := false

	// Check if the cursor is between an opening and a closing bracket
	if e.CurrentBuffer.SliceIndex > 0 && e.CurrentBuffer.SliceIndex < len(row.Chars) {
		openingBracket := row.Chars[e.CurrentBuffer.SliceIndex-1]
		cursorPos := row.Chars[e.CurrentBuffer.SliceIndex]
		if closingBracket, ok := constants.BracketPairs[rune(openingBracket)]; ok && byte(closingBracket) == cursorPos {
			isBetweenBrackets = true
		}
	}

	if e.CurrentBuffer.SliceIndex == 0 {
		newRow := config.NewRow()
		at := e.Cy
		EditorInsertRow(newRow, at, e)
	} else {
		// If we are between brackets another row will be inserted in between this row and the previous
		currentRow := e.GetCurrentRow()
		currentRow.Chars = row.Chars[:e.CurrentBuffer.SliceIndex]
		currentRow.Length = len(e.CurrentBuffer.Rows[e.Cy].Chars)

		newRow := config.Row{Chars: row.Chars[e.CurrentBuffer.SliceIndex:], IndentationLevel: currentRow.IndentationLevel}
		indentBytes := make([]byte, newRow.IndentationLevel)
		for i := 0; i < newRow.IndentationLevel; i++ {
			indentBytes[i] = byte('\t')
		}
		newRow.Chars = append(indentBytes, newRow.Chars...)

		EditorInsertRow(&newRow, e.Cy+1, e)
		e.Cx = e.LineNumberWidth
		e.CurrentBuffer.SliceIndex = 0
		if e.GetCurrentRow().IndentationLevel > 0 {
			e.CurrentBuffer.SliceIndex = constants.TAB_STOP * e.GetCurrentRow().IndentationLevel
			e.Cx = constants.TAB_STOP*e.GetCurrentRow().IndentationLevel + 5
		}
	}

	e.CurrentBuffer.NumRows++ // Update NumRows within CurrentBuffer
	e.Cy++

	// If the cursor was between brackets, insert an additional new line
	if isBetweenBrackets {
		currentRow := e.GetCurrentRow()

		newRow := config.NewRow()
		newRow.IndentationLevel = currentRow.IndentationLevel + 1
		// Create a byte slice with as many tabs as newRow.IndentationLevel
		indentBytes := make([]byte, newRow.IndentationLevel)
		for i := 0; i < newRow.IndentationLevel; i++ {
			indentBytes[i] = byte('\t')
		}

		newRow.Chars = append(indentBytes, newRow.Chars...)
		newRow.Length = len(newRow.Chars)

		EditorInsertRow(newRow, e.Cy, e)
		e.Cx = newRow.Length + e.LineNumberWidth
		e.CurrentBuffer.SliceIndex = newRow.Length
		e.CurrentBuffer.NumRows++
	}
}

func EditorInsertRow(row *config.Row, at int, e *config.Editor) {
	// Replace tabs with spaces
	convertedChars := ReplaceTabsWithSpaces(row.Chars)
	row.Chars = convertedChars
	row.Length = len(convertedChars)
	row.Idx = at // Set the index to the insertion point
	row.Highlighting = make([]byte, row.Length)
	row.Tabs = make([]byte, row.Length)

	if at < 0 || at >= len(e.CurrentBuffer.Rows) {
		// If at is outside the valid range, append the row to the end
		row.Idx = len(e.CurrentBuffer.Rows)
		highlighting.SyntaxHighlightStateMachine(row, e)
		e.CurrentBuffer.Rows = append(e.CurrentBuffer.Rows, *row)
		return
	}
	highlighting.SyntaxHighlightStateMachine(row, e)

	// Use InsertRowAtIndex to insert the row at the specified position
	e.CurrentBuffer.InsertRowAtIndex(at, *row)

	// Update the Idx of the subsequent rows
	for i := at + 1; i < len(e.CurrentBuffer.Rows); i++ {
		e.CurrentBuffer.Rows[i].Idx = i
	}
}
