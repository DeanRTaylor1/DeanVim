package actions

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func TabKeyHandler(cfg *config.EditorConfig) {
	if cfg.SliceIndex == 0 {
		cfg.GetCurrentRow().IndentationLevel++
	}
	for i := 0; i < constants.TAB_STOP; i++ {
		EditorInsertChar(' ', cfg)
	}
	MapTabs(cfg)
}

func EnterKeyHandler(cfg *config.EditorConfig) {
	action := cfg.CurrentBuffer.NewEditorAction(*cfg.GetCurrentRow().DeepCopy(), cfg.Cy, constants.ACTION_INSERT_ROW, cfg.GetCurrentRow().Length, cfg.Cx, cfg.GetCurrentRow(), func() { EditorInsertNewLine(cfg) })
	cfg.CurrentBuffer.AppendUndo(*action, cfg.UndoHistory)
	EditorInsertNewLine(cfg)
}

func QuitKeyHandler(cfg *config.EditorConfig) bool {
	if cfg.Dirty > 0 && cfg.QuitTimes > 0 {
		EditorSetStatusMessage(cfg, "WARNING!!! File has unsaved changes. Press Ctrl-Q %d more times to quit.", cfg.QuitTimes)
		cfg.QuitTimes--
		return false
	}
	fmt.Print(constants.ESCAPE_CLEAR_SCREEN)
	fmt.Print(constants.ESCAPE_MOVE_TO_HOME_POS)
	os.Exit(0)
	return true
}

func SaveKeyHandler(cfg *config.EditorConfig) {
	msg, err := EditorSave(cfg)
	if err != nil {
		EditorSetStatusMessage(cfg, "%s", err.Error())
		return
	}
	EditorSetStatusMessage(cfg, "%s", msg)
	return
}

func HomeKeyHandler(cfg *config.EditorConfig) {
	cfg.Cx = cfg.LineNumberWidth
	cfg.SliceIndex = 0
}

func EndKeyHandler(cfg *config.EditorConfig) error {
	if cfg.Cy == cfg.CurrentBuffer.NumRows {
		return errors.New("Can not go to end of this row")
	}
	cfg.Cx = cfg.CurrentBuffer.Rows[cfg.Cy].Length + cfg.LineNumberWidth
	cfg.SliceIndex = cfg.CurrentBuffer.Rows[cfg.Cy].Length
	return nil
}

func createActionForUndo(cfg *config.EditorConfig, cb func()) *config.EditorAction {
	prevRowLength := 0
	action := cfg.CurrentBuffer.NewEditorAction(*cfg.GetCurrentRow().DeepCopy(), cfg.Cy, constants.ACTION_UPDATE_ROW, prevRowLength, cfg.Cx, nil, cb)

	if cfg.Cy > 0 && cfg.SliceIndex == 0 {
		action.ActionType = constants.ACTION_APPEND_ROW_TO_PREVIOUS
		action.PrevRow = cfg.CurrentBuffer.Rows[cfg.Cy-1]
		action.Cx = cfg.LineNumberWidth
	}

	return action
}

func handleDeleteKey(cfg *config.EditorConfig, char rune) {
	if char == constants.DEL_KEY {
		EditorMoveCursor(constants.ARROW_RIGHT, cfg)
	}
}

func deleteTabOrChar(cfg *config.EditorConfig) {
	currentRow := cfg.GetCurrentRow()

	if cfg.SliceIndex > 0 && len(currentRow.Tabs) > 0 && currentRow.Tabs[cfg.SliceIndex-1] == constants.HL_TAB_KEY {
		startOfTab := cfg.SliceIndex - 1
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
		cfg.GetCurrentRow().IndentationLevel--
	} else {
		EditorDelChar(cfg)
	}
}

func DeleteHandler(cfg *config.EditorConfig, char rune) {
	action := createActionForUndo(cfg, func() { handleDeleteKey(cfg, char); deleteTabOrChar(cfg) })
	cfg.CurrentBuffer.AppendUndo(*action, cfg.UndoHistory)

	handleDeleteKey(cfg, char)
	deleteTabOrChar(cfg)
}

func PageJumpHandler(cfg *config.EditorConfig, char rune) {
	rows := cfg.ScreenRows
	for rows > 0 {
		if char == constants.PAGE_UP {
			EditorMoveCursor(constants.ARROW_UP, cfg)
		} else {
			EditorMoveCursor(constants.ARROW_DOWN, cfg)
		}
		rows--
	}
}

func IsClosingBracket(char rune) bool {
	for _, closingBracket := range constants.BracketPairs {
		if char == closingBracket {
			return true
		}
	}
	return false
}

func HandleCharInsertion(cfg *config.EditorConfig, char rune) {
	if closingBracket, ok := constants.BracketPairs[char]; ok {
		EditorInsertChar(char, cfg)
		EditorInsertChar(closingBracket, cfg)
		cfg.SliceIndex--
		cfg.Cx--
	} else {
		EditorInsertChar(char, cfg)
	}
}

func InsertCharHandler(cfg *config.EditorConfig, char rune) {
	currentRow := *cfg.GetCurrentRow().DeepCopy()
	action := cfg.CurrentBuffer.NewEditorAction(currentRow, cfg.Cy, constants.ACTION_UPDATE_ROW, 0, cfg.Cx, nil, func() { HandleCharInsertion(cfg, char) })
	cfg.CurrentBuffer.AppendUndo(*action, cfg.UndoHistory)

	HandleCharInsertion(cfg, char)
}

func ControlCHandler(buffer *bytes.Buffer, c rune, cColor int) {
	sym := '?'
	if c <= 26 {
		sym = rune(int(c) + int('@'))
	}
	buffer.WriteString("\x1b[7m")
	buffer.WriteRune(sym)
	buffer.WriteString("\x1b[m")
	if cColor != -1 {
		buffer.WriteString(fmt.Sprintf("\x1b[%dm", cColor))
	}
}

func FormatFindResultHandler(buffer *bytes.Buffer, c byte) {
	buffer.WriteString(constants.ESCAPE_HIDE_CURSOR)
	buffer.WriteString(constants.FOREGROUND_RESET)
	buffer.WriteString(constants.BACKGROUND_YELLOW)
	buffer.WriteByte(c)
	buffer.WriteString(constants.BACKGROUND_RESET)
}

func NormalFormatHandler(buffer *bytes.Buffer, c byte, cColor int) {
	if cColor != -1 {
		buffer.WriteString(constants.FOREGROUND_RESET)
		cColor = -1
	}
	buffer.WriteByte(c)
}

func ColorFormatHandler(buffer *bytes.Buffer, c byte, cColor *int, hl byte) {
	color := int(highlighting.EditorSyntaxToColor(hl))
	if color != *cColor {
		buffer.WriteString(fmt.Sprintf("\x1b[%dm", color))
		*cColor = color
	}
	buffer.WriteByte(c)
	buffer.WriteString(constants.FOREGROUND_RESET)
	*cColor = -1
}

func HideCursorIfSearching(buffer *bytes.Buffer, cfg *config.EditorConfig) {
	if cfg.CurrentBuffer.SearchState.Searching == true {
		buffer.WriteString(constants.ESCAPE_HIDE_CURSOR)
	}
}

func WriteWelcomeIfNoFile(buffer *bytes.Buffer, screenCols int, screenRows int, i int, cfg *config.EditorConfig) {
	if cfg.CurrentBuffer.NumRows == 0 && i == screenRows/3 {
		DrawWelcomeMessage(buffer, screenCols)
	} else {
		buffer.WriteByte(byte(constants.TILDE))
	}
}

func CountSpaces(cfg *config.EditorConfig, rowLength int, j int, fileRow int) (spaceCount int) {
	spaceCount = 0
	for k := j; k < j+constants.TAB_STOP; k++ {
		if k >= rowLength || cfg.CurrentBuffer.Rows[fileRow].Chars[cfg.ColOff+k] != ' ' {
			break
		}
		spaceCount++
	}
	return spaceCount
}

func AppendTabOrRowIndentBar(cfg *config.EditorConfig, j *int, buffer *bytes.Buffer, fileRow int, rowLength int) {
	nextCharIndex := *j + constants.TAB_STOP
	if nextCharIndex < rowLength && cfg.CurrentBuffer.Rows[fileRow].Chars[cfg.ColOff+nextCharIndex] != '}' {
		buffer.WriteString(strings.Repeat(" ", constants.TAB_STOP-1))
		buffer.WriteString(constants.TEXT_BRIGHT_BLACK)
		buffer.WriteString("│")
		buffer.WriteString(constants.FOREGROUND_RESET)
	} else {
		// If the next character is a '}', just append the spaces
		buffer.WriteString(strings.Repeat(" ", constants.TAB_STOP))
	}
}
