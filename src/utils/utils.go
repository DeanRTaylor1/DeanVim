package utils

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
