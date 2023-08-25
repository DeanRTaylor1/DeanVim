package utils

func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func CTRL_KEY(ch rune) rune {
	return ch & 0x1f
}
