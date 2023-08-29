package actions

import (
	"fmt"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func EditorUpdateRow(row *config.Row, cfg *config.EditorConfig) {
	if cfg.Cy < 0 {
		return
	}
	currentRow := cfg.GetCurrentRow()

	currentRow.Chars = row.Chars
	currentRow.Length = row.Length
	currentRow.Highlighting = make([]byte, row.Length)
	highlighting.Fill(currentRow.Highlighting, constants.HL_NORMAL)
	currentRow.Tabs = make([]byte, currentRow.Length)
	MapTabs(cfg)

	highlighting.SyntaxHighlightStateMachine(&cfg.CurrentBuffer.Rows[cfg.Cy], cfg)
}

func ReplaceTabsWithSpaces(line []byte) []byte {
	var result []byte
	for _, b := range line {
		if b == '\t' {
			spacesNeeded := constants.TAB_STOP - (len(result) % constants.TAB_STOP)
			for j := 0; j < spacesNeeded; j++ {
				result = append(result, byte(constants.SPACE_RUNE))
			}
		} else {
			result = append(result, b)
		}
	}
	return result
}

func MapTabs(cfg *config.EditorConfig) {
	currentRow := &cfg.CurrentBuffer.Rows[cfg.Cy]

	if len(currentRow.Tabs) != len(currentRow.Chars) {
		currentRow.Tabs = make([]byte, len(currentRow.Chars))
	}

	for i := 0; i < len(currentRow.Chars); {
		if currentRow.Chars[i] == ' ' && i+constants.TAB_STOP <= len(currentRow.Chars) {
			isTabs := true
			for j := 1; j < constants.TAB_STOP; j++ {
				if currentRow.Chars[i+j] != ' ' {
					isTabs = false
					break
				}
			}
			if isTabs {
				for j := 0; j < constants.TAB_STOP; j++ {
					currentRow.Tabs[i+j] = constants.HL_TAB_KEY
				}
				i += constants.TAB_STOP
				continue
			}
		}
		i++
	}
}

func EditorScroll(cfg *config.EditorConfig) {
	if cfg.Cy < cfg.RowOff {
		cfg.RowOff = cfg.Cy
	}
	if cfg.Cy >= cfg.RowOff+cfg.ScreenRows {
		cfg.RowOff = cfg.Cy - cfg.ScreenRows + 1
	}
	if cfg.Cx < cfg.ColOff {
		cfg.ColOff = cfg.Cx - cfg.LineNumberWidth
	}
	if cfg.Cx >= cfg.ColOff+cfg.ScreenCols {
		cfg.ColOff = cfg.Cx - cfg.ScreenCols + 1
	}
}

func SetCursorPos(x int, y int) string {
	return fmt.Sprintf(constants.ESCAPE_MOVE_TO_COORDS, x, y)
}
