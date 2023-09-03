package core

import (
	"bytes"
	"os"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func FullRefresh(e *config.Editor, buffer *bytes.Buffer) {
	buffer.WriteString(constants.ESCAPE_MOVE_TO_HOME_POS)
	buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
	if e.IsBrowsingFiles() {
		DrawFileBrowser(buffer, e, 0, e.ScreenRows)
	} else {
		EditorDrawRows(buffer, e, 0, e.ScreenRows)
	}
}

func DrawAllLineNumbers(buffer *bytes.Buffer, e *config.Editor) {
	for i := 1; i <= e.ScreenRows; i++ {
		fileRow := i + e.RowOff - 1
		if fileRow < 0 {
			continue
		}
		cursorPosition := SetCursorPos(i+1, 1)
		buffer.WriteString(cursorPosition)
		DrawLineNumbers(buffer, fileRow, e)
	}
}

func PartialRefresh(e *config.Editor, buffer *bytes.Buffer, startRow, endRow int) {
	if !e.IsBrowsingFiles() {
		DrawAllLineNumbers(buffer, e)
	}
	cursorPosition := SetCursorPos(startRow+1, 6)
	if e.IsBrowsingFiles() {
		cursorPosition = SetCursorPos(startRow+(len(e.InstructionsLines())+1), 0)
	}
	buffer.WriteString(cursorPosition)
	buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)

	if e.IsBrowsingFiles() {
		DrawFileBrowser(buffer, e, startRow, endRow)
	} else {
		EditorDrawRows(buffer, e, startRow, endRow)
	}
}

func SingleLineRefresh(e *config.Editor, buffer *bytes.Buffer, startRow, endRow int) {
	if e.IsBrowsingFiles() {
		cursorPosition := SetCursorPos(startRow+len(e.InstructionsLines())+1, 0)
		buffer.WriteString(cursorPosition)
		buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
		DrawFileBrowser(buffer, e, startRow, endRow)
	} else {
		cursorPosition := SetCursorPos(startRow+1, 0)
		buffer.WriteString(cursorPosition)
		buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
		EditorDrawRows(buffer, e, startRow, endRow)

	}
}

func EditorRefreshScreen(e *config.Editor, lastKeyPress rune) {
	if lastKeyPress == constants.NO_OP {
		return
	}
	var buffer bytes.Buffer
	EditorScroll(e)
	buffer.WriteString(constants.ESCAPE_HIDE_CURSOR)

	// Cursor position adjustment logic for non file browser modes
	if !e.IsBrowsingFiles() && e.Cx < e.LineNumberWidth {
		e.Cx = e.LineNumberWidth
	}

	if e.IsBrowsingFiles() && (e.Cy < 5 || e.Cy > len(e.FileBrowserItems)+len(e.InstructionsLines())) {
		e.Cx = 0
		e.Cy = len(e.InstructionsLines())
	}

	if !e.IsBrowsingFiles() && (e.SpecialRefreshCase()) {
		FullRefresh(e, &buffer)
	} else {
		startRow := e.Cy - 2
		endRow := e.Cy + 2
		if startRow < 0 {
			startRow = 0
		}
		switch lastKeyPress {
		case constants.INITIAL_REFRESH, constants.ENTER_KEY, constants.BACKSPACE, constants.DEL_KEY, utils.CTRL_KEY(lastKeyPress), constants.PAGE_DOWN, constants.PAGE_UP:
			FullRefresh(e, &buffer)
		case constants.ARROW_DOWN, constants.ARROW_UP, constants.ARROW_LEFT, constants.ARROW_RIGHT:
			PartialRefresh(e, &buffer, startRow, endRow)
		default:
			SingleLineRefresh(e, &buffer, 0, e.Cy)
		}
	}

	// Draw status and message bars
	statusBarPosition := SetCursorPos(e.ScreenRows+1, 0)
	buffer.WriteString(statusBarPosition)
	EditorDrawStatusBar(&buffer, e)
	EditorDrawMessageBar(&buffer, e)

	// Set cursor position
	cursorPosition := SetCursorPos((e.Cy-e.RowOff)+1, (e.Cx-e.ColOff)+1)
	buffer.WriteString(cursorPosition)

	buffer.WriteString(constants.ESCAPE_SHOW_CURSOR)

	// Write to stdout
	os.Stdout.Write(buffer.Bytes())
}
