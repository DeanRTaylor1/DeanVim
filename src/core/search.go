package core

import (
	"strings"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func EditorFindCallback(buf []rune, c rune, cfg *config.Editor, trigger bool) {
	if len(cfg.CurrentBuffer.SearchState.SavedHl) > 0 {
		sl := cfg.CurrentBuffer.SearchState.SavedHlLine
		cfg.CurrentBuffer.Rows[sl].Highlighting = cfg.CurrentBuffer.SearchState.SavedHl
	}

	if c == '\r' || c == '\x1b' {
		cfg.CurrentBuffer.SearchState.LastMatch = -1
		cfg.CurrentBuffer.SearchState.Direction = 1
		return
	} else if c == constants.ARROW_RIGHT || c == constants.ARROW_DOWN {
		cfg.CurrentBuffer.SearchState.Direction = 1
	} else if c == constants.ARROW_LEFT || c == constants.ARROW_UP {
		cfg.CurrentBuffer.SearchState.Direction = -1
	} else {
		cfg.CurrentBuffer.SearchState.LastMatch = -1
		cfg.CurrentBuffer.SearchState.Direction = 1
	}

	if cfg.CurrentBuffer.SearchState.LastMatch == -1 {
		cfg.CurrentBuffer.SearchState.Direction = 1
	}
	current := cfg.CurrentBuffer.SearchState.LastMatch
	for i := 0; i < cfg.CurrentBuffer.NumRows; i++ {
		current += cfg.CurrentBuffer.SearchState.Direction
		if current == -1 {
			current = cfg.CurrentBuffer.NumRows - 1
		} else if current == cfg.CurrentBuffer.NumRows {
			current = 0
		}

		row := cfg.CurrentBuffer.Rows[current].Chars
		matchIndex := strings.Index(string(row), string(buf))
		if matchIndex != -1 {
			cfg.CurrentBuffer.SearchState.LastMatch = current
			cfg.Cy = current
			cfg.Cx = matchIndex + cfg.LineNumberWidth
			cfg.CurrentBuffer.SliceIndex = matchIndex
			cfg.RowOff = cfg.CurrentBuffer.NumRows

			cfg.CurrentBuffer.SearchState.SavedHlLine = current
			cfg.CurrentBuffer.SearchState.SavedHl = make([]byte, len(cfg.CurrentBuffer.Rows[cfg.Cy].Highlighting))
			copy(cfg.CurrentBuffer.SearchState.SavedHl, cfg.CurrentBuffer.Rows[cfg.Cy].Highlighting)

			for i := 0; i < len(buf); i++ {
				cfg.CurrentBuffer.Rows[cfg.Cy].Highlighting[matchIndex+i] = constants.HL_MATCH
			}
			break
		}
	}
}

func EditorFind(cfg *config.Editor) {
	cfg.CurrentBuffer.SearchState.Searching = true
	cx := cfg.Cx
	sliceIndex := cfg.CurrentBuffer.SliceIndex
	cy := cfg.Cy
	rowOff := cfg.RowOff
	colOff := cfg.ColOff

	query := EditorPrompt("Search: (ESC to cancel)", EditorFindCallback, cfg)

	if query == nil {
		cfg.Cx = cx
		cfg.CurrentBuffer.SliceIndex = sliceIndex
		cfg.Cy = cy
		cfg.RowOff = rowOff
		cfg.ColOff = colOff
	}
	cfg.CurrentBuffer.SearchState.Searching = false
}
