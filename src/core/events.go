package core

import (
	"bufio"
	"fmt"
	"time"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func EventHandlerMain(reader *bufio.Reader, e *config.Editor) rune {
	char, err := ReadKey(reader)
	if err != nil {
		panic(err)
	}

	if e.EditorMode == constants.EDITOR_MODE_NORMAL {
		char = NormalModeEventsHandler(char, e)
	} else if e.EditorMode == constants.EDITOR_MODE_INSERT {
		InsertModeEventsHandler(char, e)
	} else if e.IsBrowsingFiles() {
		char = FileBrowserEventsHandler(char, e)
	} else if e.EditorMode == constants.EDITOR_MODE_VISUAL {
		char = VisualModeEventsHandler(char, e)
	}
	return char
}

func FileBrowserCursorMovements(key rune, e *config.Editor) {
	switch key {
	case rune(constants.ARROW_LEFT):
		if e.Cx <= 0 {
			return
		}
		e.MoveCursorLeft()
	case rune(constants.ARROW_RIGHT):
		if e.Cy == len(e.FileBrowserItems)+len(e.InstructionsLines()) {
			break
		}
		e.MoveCursorRight()
	case rune(constants.ARROW_DOWN):
		if e.Cy < len(e.FileBrowserItems)+len(e.InstructionsLines()) {
			e.MoveCursorDown()
		}
	case rune(constants.ARROW_UP):
		if e.Cy >= 5 {
			e.MoveCursorUp()
		}
	}
}

func EditorCursorMovements(key rune, e *config.Editor) {
	var row []byte = []byte{}
	if e.Cy < e.CurrentBuffer.NumRows {
		row = e.CurrentBuffer.Rows[e.Cy].Chars
	}
	// spacesNeeded := TAB_STOP - (e.Cx % TAB_STOP)
	switch key {
	case rune(constants.ARROW_LEFT):
		if e.CurrentBuffer.SliceIndex != 0 {
			e.MoveCursorLeft()
		} else if e.Cy > 0 && e.Cy < e.CurrentBuffer.NumRows {
			e.MoveCursorUp()
			e.Cx = (e.GetCurrentRow().Length) + e.LineNumberWidth
			e.CurrentBuffer.SliceIndex = e.GetCurrentRow().Length
		}
	case rune(constants.ARROW_RIGHT):
		if e.Cy == e.CurrentBuffer.NumRows {
			break
		}
		if e.CurrentBuffer.SliceIndex < (e.GetCurrentRow().Length) {
			e.MoveCursorRight()
		} else if e.Cx-e.LineNumberWidth >= e.GetCurrentRow().Length && e.Cy < len(e.CurrentBuffer.Rows)-1 {
			e.MoveCursorDown()
			e.Cx = e.LineNumberWidth
			e.CurrentBuffer.SliceIndex = 0
		}
	case rune(constants.ARROW_DOWN):
		if e.Cy < e.CurrentBuffer.NumRows {
			e.MoveCursorDown()
		}
	case rune(constants.ARROW_UP):
		if e.Cy != 0 {
			e.MoveCursorUp()
		}
	}

	if e.Cy < e.CurrentBuffer.NumRows {
		row = e.CurrentBuffer.Rows[e.Cy].Chars
	} else {
		row = []byte{}
	}

	rowLen := len(row)
	if e.CurrentBuffer.SliceIndex > rowLen {
		e.Cx = rowLen + e.LineNumberWidth
		e.CurrentBuffer.SliceIndex = rowLen
	}
}

func EditorMoveCursor(key rune, e *config.Editor) {
	if e.IsBrowsingFiles() {
		FileBrowserCursorMovements(key, e)
	} else {
		EditorCursorMovements(key, e)
	}
}

func EditorSetStatusMessage(e *config.Editor, format string, a ...interface{}) {
	e.StatusMsg = fmt.Sprintf(format, a...)
	e.StatusMsgTime = time.Now()
}
