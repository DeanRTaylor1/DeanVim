package actions

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func editorRowInsertChar(row *config.Row, at int, char rune, cfg *config.EditorConfig) {
	row.Chars = append(row.Chars, 0)
	copy(row.Chars[at+1:], row.Chars[at:])
	row.Chars[at] = byte(char)

	row.Length = len(row.Chars)
	row.Highlighting = make([]byte, row.Length)
	highlighting.Fill(row.Highlighting, constants.HL_NORMAL)
	row.Tabs = make([]byte, row.Length)

	highlighting.SyntaxHighlightStateMachine(row, cfg)

	cfg.Dirty++
}

func EditorInsertChar(char rune, cfg *config.EditorConfig) {
	if cfg.Cy == cfg.CurrentBuffer.NumRows {
		EditorInsertRow(config.NewRow(), -1, cfg)
		cfg.CurrentBuffer.NumRows++
	}
	editorRowInsertChar(&cfg.CurrentBuffer.Rows[cfg.Cy], cfg.SliceIndex, char, cfg)

	cfg.MoveCursorRight()
}

func EditorInsertNewLine(cfg *config.EditorConfig) {
	row := cfg.CurrentBuffer.Rows[cfg.Cy]
	isBetweenBrackets := false

	// Check if the cursor is between an opening and a closing bracket
	if cfg.SliceIndex > 0 && cfg.SliceIndex < len(row.Chars) {
		openingBracket := row.Chars[cfg.SliceIndex-1]
		cursorPos := row.Chars[cfg.SliceIndex]
		if closingBracket, ok := constants.BracketPairs[rune(openingBracket)]; ok && byte(closingBracket) == cursorPos {
			isBetweenBrackets = true
		}
	}

	if cfg.SliceIndex == 0 {
		newRow := config.NewRow()
		at := cfg.Cy
		EditorInsertRow(newRow, at, cfg)
	} else {
		// If we are between brackets another row will be inserted in between this row and the previous
		currentRow := cfg.GetCurrentRow()
		currentRow.Chars = row.Chars[:cfg.SliceIndex]
		currentRow.Length = len(cfg.CurrentBuffer.Rows[cfg.Cy].Chars)

		newRow := config.Row{Chars: row.Chars[cfg.SliceIndex:], IndentationLevel: currentRow.IndentationLevel}
		indentBytes := make([]byte, newRow.IndentationLevel)
		for i := 0; i < newRow.IndentationLevel; i++ {
			indentBytes[i] = byte('\t')
		}
		newRow.Chars = append(indentBytes, newRow.Chars...)

		EditorInsertRow(&newRow, cfg.Cy+1, cfg)
		cfg.Cx = cfg.LineNumberWidth
		cfg.SliceIndex = 0
		if cfg.GetCurrentRow().IndentationLevel > 0 {
			cfg.SliceIndex = constants.TAB_STOP * cfg.GetCurrentRow().IndentationLevel
			cfg.Cx = constants.TAB_STOP*cfg.GetCurrentRow().IndentationLevel + 5
		}
	}

	cfg.CurrentBuffer.NumRows++ // Update NumRows within CurrentBuffer
	cfg.Cy++

	// If the cursor was between brackets, insert an additional new line
	if isBetweenBrackets {
		currentRow := cfg.GetCurrentRow()

		newRow := config.NewRow()
		newRow.IndentationLevel = currentRow.IndentationLevel + 1
		// Create a byte slice with as many tabs as newRow.IndentationLevel
		indentBytes := make([]byte, newRow.IndentationLevel)
		for i := 0; i < newRow.IndentationLevel; i++ {
			indentBytes[i] = byte('\t')
		}

		newRow.Chars = append(indentBytes, newRow.Chars...)
		newRow.Length = len(newRow.Chars)

		EditorInsertRow(newRow, cfg.Cy, cfg)
		cfg.Cx = newRow.Length + cfg.LineNumberWidth
		cfg.SliceIndex = newRow.Length
		cfg.CurrentBuffer.NumRows++
	}
}

func EditorInsertRow(row *config.Row, at int, cfg *config.EditorConfig) {
	// Replace tabs with spaces
	convertedChars := ReplaceTabsWithSpaces(row.Chars)
	row.Chars = convertedChars
	row.Length = len(convertedChars)
	row.Idx = at // Set the index to the insertion point
	row.Highlighting = make([]byte, row.Length)
	row.Tabs = make([]byte, row.Length)

	if at < 0 || at >= len(cfg.CurrentBuffer.Rows) {
		// If at is outside the valid range, append the row to the end
		row.Idx = len(cfg.CurrentBuffer.Rows)
		highlighting.SyntaxHighlightStateMachine(row, cfg)
		cfg.CurrentBuffer.Rows = append(cfg.CurrentBuffer.Rows, *row)
		return
	}
	highlighting.SyntaxHighlightStateMachine(row, cfg)

	// Use InsertRowAtIndex to insert the row at the specified position
	cfg.CurrentBuffer.InsertRowAtIndex(at, *row)

	// Update the Idx of the subsequent rows
	for i := at + 1; i < len(cfg.CurrentBuffer.Rows); i++ {
		cfg.CurrentBuffer.Rows[i].Idx = i
	}
}
