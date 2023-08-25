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
	editorRowInsertChar(&cfg.CurrentBuffer.Rows[cfg.Cy], cfg.Cx, char, cfg)

	cfg.Cx++
}

func EditorInsertNewLine(cfg *config.EditorConfig) {
	row := cfg.CurrentBuffer.Rows[cfg.Cy]
	isBetweenBrackets := false

	// Check if the cursor is between an opening and a closing bracket
	if cfg.Cx > 0 && cfg.Cx < len(row.Chars) {
		openingBracket := row.Chars[cfg.Cx-1]
		cursorPos := row.Chars[cfg.Cx]
		if closingBracket, ok := constants.BracketPairs[rune(openingBracket)]; ok && byte(closingBracket) == cursorPos {
			isBetweenBrackets = true
		}
	}

	if cfg.Cx == 0 {
		newRow := config.NewRow()
		at := cfg.Cy
		EditorInsertRow(newRow, at, cfg)
	} else {
		cfg.CurrentBuffer.Rows[cfg.Cy].Chars = row.Chars[:cfg.Cx]
		cfg.CurrentBuffer.Rows[cfg.Cy].Length = len(cfg.CurrentBuffer.Rows[cfg.Cy].Chars)
		newRow := config.Row{Chars: row.Chars[cfg.Cx:]}
		EditorInsertRow(&newRow, cfg.Cy+1, cfg)
		cfg.Cx = 0
	}

	cfg.CurrentBuffer.NumRows++ // Update NumRows within CurrentBuffer
	cfg.Cy++

	// If the cursor was between brackets, insert an additional new line
	if isBetweenBrackets {
		newRow := config.NewRow()
		EditorInsertRow(newRow, cfg.Cy, cfg)
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
	highlighting.SyntaxHighlightStateMachine(row, cfg)

	if at < 0 || at >= len(cfg.CurrentBuffer.Rows) {
		// If at is outside the valid range, append the row to the end
		cfg.CurrentBuffer.Rows = append(cfg.CurrentBuffer.Rows, *row)
		return
	}

	// If at is within the valid range, insert the row at the specified position
	cfg.CurrentBuffer.Rows = append(cfg.CurrentBuffer.Rows[:at], append([]config.Row{*row}, cfg.CurrentBuffer.Rows[at:]...)...)

	// Update the Idx of the subsequent rows
	for i := at + 1; i < len(cfg.CurrentBuffer.Rows); i++ {
		cfg.CurrentBuffer.Rows[i].Idx = i
	}
}
