package highlighting

import (
	"strings"
	"unicode"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func isSeparator(c rune) bool {
	return unicode.IsSpace(c) || c == '\x00' || strings.ContainsRune(",.()+-/*=~%<>[];", c)
}

func EditorUpdateSyntax(row *config.Row) {
	row.Highlighting = make([]byte, row.Length)
	prevSep := true
	for i := 0; i < row.Length; i++ {
		c := row.Chars[i]
		prevHl := byte(constants.HL_NORMAL)
		if i > 0 {
			prevHl = row.Highlighting[i-1]
		}
		if unicode.IsDigit(rune(c)) && (prevSep || prevHl == constants.HL_NUMBER) || (c == '.' && prevHl == constants.HL_NUMBER) {
			row.Highlighting[i] = constants.HL_NUMBER
			prevSep = false
			continue
		}
		prevSep = isSeparator(rune(c))
	}
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
