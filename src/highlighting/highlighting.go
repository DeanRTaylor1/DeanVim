package highlighting

import (
	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func EditorUpdateSyntax(row *config.Row) {
	hl := make([]byte, row.Length)
	for i, c := range row.Chars {
		hl[i] = constants.HL_NORMAL
		if utils.IsDigit(c) {
			hl[i] = constants.HL_NUMBER
		}
	}
	row.Highlighting = hl
}

func EditorSyntaxToColor(highlight byte) byte {
	switch highlight {
	case constants.HL_NUMBER:
		return 31
	case constants.HL_MATCH:
		return 34
	default:
		return 37
	}
}
