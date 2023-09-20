package core

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func TabKeyHandler(e *config.Editor) {
	if e.CurrentBuffer.SliceIndex == 0 {
		e.GetCurrentRow().IndentationLevel++
	}
	for i := 0; i < constants.TAB_STOP; i++ {
		EditorInsertChar(' ', e)
	}
	MapTabs(e)
}

func EnterKeyHandler(e *config.Editor) {
	action := e.CurrentBuffer.NewEditorAction(*e.GetCurrentRow().DeepCopy(), e.Cy, constants.ACTION_INSERT_ROW, e.GetCurrentRow().Length, e.Cx, e.GetCurrentRow(), func() { EditorInsertNewLine(e) })
	e.CurrentBuffer.AppendUndo(*action, e.UndoHistory)
	EditorInsertNewLine(e)
}

func QuitKeyHandler(e *config.Editor) bool {
	if e.CurrentBuffer.Dirty > 0 && e.QuitTimes > 0 {
		EditorSetStatusMessage(e, "WARNING!!! File has unsaved changes. Press Ctrl-Q %d more times to quit.", e.QuitTimes)
		e.QuitTimes--
		return false
	}

	e.RemoveBuffer(e.CurrentBuffer.Name) // Remove the current buffer

	if len(e.Buffers) > 0 {
		fullPath := filepath.Join(e.RootDirectory, e.Buffers[0].Name)
		// Load the next buffer if there's any remaining
		ReadHandler(e, fullPath)
		return false
	}

	// Clear screen and exit if no more buffers
	fmt.Print(constants.ESCAPE_CLEAR_SCREEN)
	fmt.Print(constants.ESCAPE_MOVE_TO_HOME_POS)
	os.Exit(0)
	return true
}

func SaveKeyHandler(e *config.Editor) {
	msg, err := EditorSave(e)
	if err != nil {
		EditorSetStatusMessage(e, "%s", err.Error())
		return
	}
	EditorSetStatusMessage(e, "%s", msg)
	return
}

func HomeKeyHandler(e *config.Editor) {
	e.Cx = e.LineNumberWidth
	e.CurrentBuffer.SliceIndex = 0
}

func EndKeyHandler(e *config.Editor) error {
	if e.Cy == e.CurrentBuffer.NumRows {
		return errors.New("Can not go to end of this row")
	}
	e.Cx = e.CurrentBuffer.Rows[e.Cy].Length + e.LineNumberWidth
	e.CurrentBuffer.SliceIndex = e.CurrentBuffer.Rows[e.Cy].Length
	return nil
}

func createActionForUndo(e *config.Editor, cb func()) *config.EditorAction {
	prevRowLength := 0
	action := e.CurrentBuffer.NewEditorAction(*e.GetCurrentRow().DeepCopy(), e.Cy, constants.ACTION_UPDATE_ROW, prevRowLength, e.Cx, nil, cb)

	if e.Cy > 0 && e.CurrentBuffer.SliceIndex == 0 {
		action.ActionType = constants.ACTION_APPEND_ROW_TO_PREVIOUS
		action.PrevRow = e.CurrentBuffer.Rows[e.Cy-1]
		action.Cx = e.LineNumberWidth
	}

	return action
}

func handleDeleteKey(e *config.Editor, char rune) {
	if char == constants.DEL_KEY {
		EditorMoveCursor(constants.ARROW_RIGHT, e)
	}
}

func deleteTabOrChar(e *config.Editor) {
	if e.ModalOpen {
	} else {

		currentRow := e.GetCurrentRow()

		if e.CurrentBuffer.SliceIndex > 0 && len(currentRow.Tabs) > 0 && currentRow.Tabs[e.CurrentBuffer.SliceIndex-1] == constants.HL_TAB_KEY {
			startOfTab := e.CurrentBuffer.SliceIndex - 1
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
				EditorDelChar(e)
			}
			e.GetCurrentRow().IndentationLevel--
		} else {
			EditorDelChar(e)
		}
	}
}

func DeleteHandler(e *config.Editor, char rune) {
	if e.ModalOpen {
		handleDeleteKey(e, char)
		ModalSearchDelChar(e)
	} else {
		action := createActionForUndo(e, func() { handleDeleteKey(e, char); deleteTabOrChar(e) })
		e.CurrentBuffer.AppendUndo(*action, e.UndoHistory)

		handleDeleteKey(e, char)
		deleteTabOrChar(e)

	}
}

func PageJumpHandler(e *config.Editor, char rune) {
	rows := e.ScreenRows
	for rows > 0 {
		if char == constants.PAGE_UP {
			EditorMoveCursor(constants.ARROW_UP, e)
		} else {
			EditorMoveCursor(constants.ARROW_DOWN, e)
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

func HandleCharInsertion(e *config.Editor, char rune) {
	if closingBracket, ok := constants.BracketPairs[char]; ok {
		EditorInsertChar(char, e)
		EditorInsertChar(closingBracket, e)
		e.CurrentBuffer.SliceIndex--
		e.Cx--
	} else {
		EditorInsertChar(char, e)
	}
}

func InsertCharHandler(e *config.Editor, char rune) {
	var currentRow config.Row
	var action *config.EditorAction
	if e.Cy != e.CurrentBuffer.NumRows {
		currentRow = *e.GetCurrentRow().DeepCopy()
		action = e.CurrentBuffer.NewEditorAction(currentRow, e.Cy, constants.ACTION_UPDATE_ROW, 0, e.Cx, nil, func() { HandleCharInsertion(e, char) })
	} else {
		currentRow = *config.NewRow()
		action = e.CurrentBuffer.NewEditorAction(currentRow, e.CurrentBuffer.NumRows, constants.ACTION_INSERT_CHAR_AT_EOF, 0, e.Cx, nil, func() { HandleCharInsertion(e, char) })
	}
	e.CurrentBuffer.AppendUndo(*action, e.UndoHistory)

	HandleCharInsertion(e, char)
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

func FormatSelectedTextHandler(buffer *bytes.Buffer, c byte, cColor *int, hl byte) {
	buffer.WriteString(constants.BACKGROUND_BRIGHT_BLACK)
	color := int(highlighting.EditorSyntaxToColor(hl))
	if color != *cColor {
		buffer.WriteString(fmt.Sprintf("\x1b[%dm", color))
		*cColor = color
	}
	buffer.WriteByte(c)
	buffer.WriteString(constants.BACKGROUND_RESET)
	*cColor = -1
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

func HideCursorIf(buffer *bytes.Buffer, propertyTrigger bool) {
	if propertyTrigger == true {
		buffer.WriteString(constants.ESCAPE_HIDE_CURSOR)
	}
}

func HideCursorIfSearching(buffer *bytes.Buffer, e *config.Editor) {
	if e.CurrentBuffer.SearchState.Searching == true {
		buffer.WriteString(constants.ESCAPE_HIDE_CURSOR)
	}
}

func WriteWelcomeIfNoFile(buffer *bytes.Buffer, screenCols int, screenRows int, i int, e *config.Editor) {
	if e.CurrentBuffer.NumRows == 0 && i == screenRows/3 {
		DrawWelcomeMessage(buffer, screenCols)
	} else {
		buffer.WriteByte(byte(constants.TILDE))
	}
}

func CountSpaces(e *config.Editor, rowLength int, j int, fileRow int) (spaceCount int) {
	spaceCount = 0
	for k := j; k < j+constants.TAB_STOP; k++ {
		if k >= rowLength || e.CurrentBuffer.Rows[fileRow].Chars[e.ColOff+k] != ' ' {
			break
		}
		spaceCount++
	}
	return spaceCount
}

func AppendTabOrRowIndentBar(e *config.Editor, j *int, buffer *bytes.Buffer, fileRow int, rowLength int) {
	nextCharIndex := *j + constants.TAB_STOP
	if nextCharIndex < rowLength && e.CurrentBuffer.Rows[fileRow].Chars[e.ColOff+nextCharIndex] != '}' {
		buffer.WriteString(strings.Repeat(" ", constants.TAB_STOP-1))
		buffer.WriteString(constants.TEXT_BRIGHT_BLACK)
		buffer.WriteString("│")
		buffer.WriteString(constants.FOREGROUND_RESET)
	} else {
		// If the next character is a '}', just append the spaces
		buffer.WriteString(strings.Repeat(" ", constants.TAB_STOP))
	}
}
