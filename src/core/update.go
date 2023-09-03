package core

import (
	"fmt"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func EditorUpdateRow(row *config.Row, e *config.Editor) {
	if e.Cy < 0 {
		return
	}
	currentRow := e.GetCurrentRow()

	currentRow.Chars = row.Chars
	currentRow.Length = row.Length
	currentRow.Highlighting = make([]byte, row.Length)
	highlighting.Fill(currentRow.Highlighting, constants.HL_NORMAL)
	currentRow.Tabs = make([]byte, currentRow.Length)
	MapTabs(e)

	highlighting.SyntaxHighlightStateMachine(&e.CurrentBuffer.Rows[e.Cy], e)
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

func MapTabs(e *config.Editor) {
	currentRow := &e.CurrentBuffer.Rows[e.Cy]

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

func EditorScroll(e *config.Editor) {
	if e.Cy <= e.RowOff {
		e.RowOff = e.Cy
	}
	if e.Cy >= e.RowOff+e.ScreenRows {
		e.RowOff = e.Cy - e.ScreenRows + 1
	}
	if e.Cx >= 5 && e.Cx-e.LineNumberWidth < e.ColOff {
		e.ColOff = e.Cx - e.LineNumberWidth
	}
	if e.Cx >= e.ColOff+e.ScreenCols {
		e.ColOff = e.Cx - e.ScreenCols + 1
	}
}

func SetCursorPos(x int, y int) string {
	return fmt.Sprintf(constants.ESCAPE_MOVE_TO_COORDS, x, y)
}
