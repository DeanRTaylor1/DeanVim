package utils

import (
	"strings"

	"github.com/deanrtaylor1/go-editor/constants"
)

func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func CTRL_KEY(ch rune) rune {
	return ch & 0x1f
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IsValidStartingChar(char rune, mode int) bool {
	var validChars string
	switch mode {
	case constants.EDITOR_MODE_NORMAL:
		validChars = " wb0$^GgvoOyppV" // Exclude hjkl
	case constants.EDITOR_MODE_VISUAL:
		validChars = "123456789"
	default:
		return false
	}
	return strings.ContainsRune(validChars, char)
}
