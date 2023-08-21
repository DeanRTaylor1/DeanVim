package highlighting

import (
	"unicode"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func EditorUpdateSyntax(r []byte, rowIndex int, cfg *config.EditorConfig) {
	hl := make([]byte, len(r))
	for i, c := range r {
		hl[i] = constants.HL_NORMAL
		if unicode.IsNumber(rune(c)) {
			hl[i] = constants.HL_NUMBER
		}
	}
	cfg.Highlighting[rowIndex] = hl
}

func EditorSyntaxToColor(highlight byte) byte {
	switch highlight {
	case constants.HL_NUMBER:
		return 31
	default:
		return 37
	}
}
