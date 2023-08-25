package actions

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func ProcessKeyPress(reader *bufio.Reader, cfg *config.EditorConfig) {
	char, err := ReadKey(reader)
	if err != nil {
		panic(err)
	}
	switch char {
	case constants.TAB_KEY:
		for i := 0; i < constants.TAB_STOP; i++ {
			EditorInsertChar(' ', cfg)
		}
		MapTabs(cfg)
		break
	case constants.ENTER_KEY:
		EditorInsertNewLine(cfg)
		break
	case utils.CTRL_KEY(constants.QUIT_KEY):
		if cfg.Dirty > 0 && cfg.QuitTimes > 0 {
			EditorSetStatusMessage(cfg, "WARNING!!! File has unsaved changes. Press Ctrl-Q %d more times to quit.", cfg.QuitTimes)
			cfg.QuitTimes--
			return
		}
		fmt.Print(constants.ESCAPE_CLEAR_SCREEN)
		fmt.Print(constants.ESCAPE_MOVE_TO_HOME_POS)
		os.Exit(0)
		break
	case utils.CTRL_KEY(constants.SAVE_KEY):
		msg, err := EditorSave(cfg)
		if err != nil {
			EditorSetStatusMessage(cfg, "%s", err.Error())
			break
		}
		EditorSetStatusMessage(cfg, "%s", msg)
		break
	case constants.HOME_KEY:
		cfg.Cx = 0
		break
	case constants.END_KEY:
		if cfg.Cy == cfg.CurrentBuffer.NumRows {
			break
		}
		cfg.Cx = cfg.CurrentBuffer.Rows[cfg.Cy].Length
		break
	case utils.CTRL_KEY('f'):
		EditorFind(cfg)
		break
	case constants.BACKSPACE, utils.CTRL_KEY('h'), constants.DEL_KEY:
		if char == constants.DEL_KEY {
			EditorMoveCursor(constants.ARROW_RIGHT, cfg)
		}

		currentRow := &cfg.CurrentBuffer.Rows[cfg.Cy]
		if cfg.Cx > 0 && currentRow.Tabs[cfg.Cx-1] == constants.HL_TAB_KEY {
			startOfTab := cfg.Cx - 1
			endOfTab := startOfTab
			i := 1
			for startOfTab > 0 && currentRow.Tabs[startOfTab-1] == constants.HL_TAB_KEY {
				startOfTab--
				i++
				if i == constants.TAB_STOP {
					break // Stop after finding one complete tab
				}
			}

			// Delete the entire tab
			for j := endOfTab; j >= startOfTab; j-- {
				EditorDelChar(cfg)
			}
		} else {
			EditorDelChar(cfg)
		}
		break
	case constants.PAGE_DOWN, constants.PAGE_UP:
		times := cfg.ScreenRows
		for times > 0 {
			if char == constants.PAGE_UP {
				EditorMoveCursor(constants.ARROW_UP, cfg)
			} else {
				EditorMoveCursor(constants.ARROW_DOWN, cfg)
			}
			times--
		}
	case utils.CTRL_KEY('l'), '\x1b':
		break
	case rune(constants.ARROW_DOWN), rune(constants.ARROW_UP), rune(constants.ARROW_RIGHT), rune(constants.ARROW_LEFT):
		EditorMoveCursor(char, cfg)
	default:
		if closingBracket, ok := constants.BracketPairs[char]; ok {
			EditorInsertChar(char, cfg)
			EditorInsertChar(closingBracket, cfg)
			cfg.Cx--
			break
		} else {
			EditorInsertChar(char, cfg)
			break

		}
	}
	cfg.QuitTimes = constants.QUIT_TIMES
}

func EditorMoveCursor(key rune, cfg *config.EditorConfig) {
	row := []byte{}
	if cfg.Cy < cfg.CurrentBuffer.NumRows {
		row = cfg.CurrentBuffer.Rows[cfg.Cy].Chars
	}
	// spacesNeeded := TAB_STOP - (cfg.Cx % TAB_STOP)
	switch key {
	case rune(constants.ARROW_LEFT):
		if cfg.Cx != 0 {
			cfg.Cx--
		} else if cfg.Cy > 0 {
			cfg.Cy--
			if cfg.Cy < len(cfg.CurrentBuffer.Rows) {
				cfg.Cx = (cfg.CurrentBuffer.Rows[cfg.Cy].Length)
			}
		}
		break
	case rune(constants.SAVE_KEY):
		if cfg.Cy == cfg.CurrentBuffer.NumRows {
			break
		}
		if cfg.Cx < (cfg.CurrentBuffer.Rows[cfg.Cy].Length)-1 {
			cfg.Cx++
		} else if cfg.Cx == (cfg.CurrentBuffer.Rows[cfg.Cy].Length) && cfg.Cy < len(cfg.CurrentBuffer.Rows)-1 {
			cfg.Cy++
			cfg.Cx = 0
		}
		break
	case rune(constants.ARROW_RIGHT):
		if cfg.Cy == cfg.CurrentBuffer.NumRows {
			break
		}
		if cfg.Cx < (cfg.CurrentBuffer.Rows[cfg.Cy].Length) {
			cfg.Cx++
		} else if cfg.Cx == cfg.CurrentBuffer.Rows[cfg.Cy].Length && cfg.Cy < len(cfg.CurrentBuffer.Rows)-1 {
			cfg.Cy++
			cfg.Cx = 0
		}
		break

	case rune(constants.ARROW_DOWN):
		if cfg.Cy < cfg.CurrentBuffer.NumRows {
			cfg.Cy++
		}
		break
	case rune(constants.ARROW_UP):
		if cfg.Cy != 0 {
			cfg.Cy--
		}
		break
	}

	if cfg.Cy < cfg.CurrentBuffer.NumRows {
		row = cfg.CurrentBuffer.Rows[cfg.Cy].Chars
	} else {
		row = []byte{}
	}

	rowLen := len(row)
	if cfg.Cx > rowLen {
		cfg.Cx = rowLen
	}
}

func EditorSetStatusMessage(cfg *config.EditorConfig, format string, a ...interface{}) {
	cfg.StatusMsg = fmt.Sprintf(format, a...)
	cfg.StatusMsgTime = time.Now()
}
