package actions

import (
	"bufio"

	"github.com/deanrtaylor1/go-editor/constants"
)

func ReadKey(reader *bufio.Reader) (rune, error) {
	char, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	if char == '\t' {
		return constants.TAB_KEY, nil
	}
	// readRune returns one byte so we check if that byte is an escape character
	// That means an arrow key could have been pressed so we replace that arrow with our keys for navigation
	// which copy vims mappings

	if char == '\x1b' {
		seq := make([]rune, 3)
		seq[0], _, err = reader.ReadRune()
		if err != nil {
			return '\x1b', nil
		}
		seq[1], _, err = reader.ReadRune()
		if err != nil {
			return '\x1b', nil
		}

		if seq[0] == '[' {
			if seq[1] >= '0' && seq[1] <= '9' {
				seq[2], _, err = reader.ReadRune()
				if err != nil {
					return '\x1b', nil
				}
				if seq[2] == '~' {
					switch seq[1] {
					case '1':
						return constants.HOME_KEY, nil
					case '3':
						return constants.DEL_KEY, nil
					case '4':
						return constants.END_KEY, nil
					case '5':
						return constants.PAGE_UP, nil
					case '6':
						return constants.PAGE_DOWN, nil
					case '7':
						return constants.HOME_KEY, nil
					case '8':
						return constants.END_KEY, nil
					}
				}
			} else {
				switch seq[1] {
				case 'A':
					return constants.ARROW_UP, nil // Up
				case 'B':
					return constants.ARROW_DOWN, nil // Down
				case 'C':
					return constants.ARROW_RIGHT, nil // Right
				case 'D':
					return constants.ARROW_LEFT, nil // Left
				case 'H':
					return constants.HOME_KEY, nil
				case 'F':
					return constants.END_KEY, nil
				}
			}
		} else if seq[0] == 'O' {
			switch seq[1] {
			case 'H':
				return constants.HOME_KEY, nil
			case 'F':
				return constants.END_KEY, nil
			}
		}
		return '\x1b', nil
	} else {
		return char, nil // other keypresses
	}
}
