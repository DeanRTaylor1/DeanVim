package core

import (
	"strings"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func EditorFindCallback(buf []rune, c rune, e *config.Editor, trigger bool) {
	if len(e.CurrentBuffer.SearchState.SavedHl) > 0 {
		sl := e.CurrentBuffer.SearchState.SavedHlLine
		e.CurrentBuffer.Rows[sl].Highlighting = e.CurrentBuffer.SearchState.SavedHl
	}

	if c == '\r' || c == '\x1b' {
		e.CurrentBuffer.SearchState.LastMatch = -1
		e.CurrentBuffer.SearchState.Direction = 1
		return
	} else if c == constants.ARROW_RIGHT || c == constants.ARROW_DOWN {
		e.CurrentBuffer.SearchState.Direction = 1
	} else if c == constants.ARROW_LEFT || c == constants.ARROW_UP {
		e.CurrentBuffer.SearchState.Direction = -1
	} else {
		e.CurrentBuffer.SearchState.LastMatch = -1
		e.CurrentBuffer.SearchState.Direction = 1
	}

	if e.CurrentBuffer.SearchState.LastMatch == -1 {
		e.CurrentBuffer.SearchState.Direction = 1
	}
	current := e.CurrentBuffer.SearchState.LastMatch
	for i := 0; i < e.CurrentBuffer.NumRows; i++ {
		current += e.CurrentBuffer.SearchState.Direction
		if current == -1 {
			current = e.CurrentBuffer.NumRows - 1
		} else if current == e.CurrentBuffer.NumRows {
			current = 0
		}

		row := e.CurrentBuffer.Rows[current].Chars
		matchIndex := strings.Index(string(row), string(buf))
		if matchIndex != -1 {
			e.CurrentBuffer.SearchState.LastMatch = current
			e.Cy = current
			e.Cx = matchIndex + e.LineNumberWidth
			e.CurrentBuffer.SliceIndex = matchIndex
			e.RowOff = e.CurrentBuffer.NumRows

			e.CurrentBuffer.SearchState.SavedHlLine = current
			e.CurrentBuffer.SearchState.SavedHl = make([]byte, len(e.CurrentBuffer.Rows[e.Cy].Highlighting))
			copy(e.CurrentBuffer.SearchState.SavedHl, e.CurrentBuffer.Rows[e.Cy].Highlighting)

			for i := 0; i < len(buf); i++ {
				e.CurrentBuffer.Rows[e.Cy].Highlighting[matchIndex+i] = constants.HL_MATCH
			}
			break
		}
	}
}

func EditorFind(e *config.Editor) {
	e.CurrentBuffer.SearchState.Searching = true
	cx := e.Cx
	sliceIndex := e.CurrentBuffer.SliceIndex
	cy := e.Cy
	rowOff := e.RowOff
	colOff := e.ColOff

	query := EditorPrompt("Search: (ESC to cancel)", EditorFindCallback, e)

	if query == nil {
		e.Cx = cx
		e.CurrentBuffer.SliceIndex = sliceIndex
		e.Cy = cy
		e.RowOff = rowOff
		e.ColOff = colOff
	}
	e.CurrentBuffer.SearchState.Searching = false
}
