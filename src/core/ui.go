package core

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"time"
	"unicode"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func DrawLineNumbers(buffer *bytes.Buffer, fileRow int, e *config.Editor) {
	relativeLineNumber := int(math.Abs(float64(e.Cy - fileRow)))
	lineNumber := fmt.Sprintf("%4d ", relativeLineNumber)

	if fileRow == e.Cy {
		buffer.WriteString(constants.TEXT_BRIGHT_WHITE)
		lineNumber = fmt.Sprintf("~%3d ", fileRow+1)
	} else {
		buffer.WriteString(constants.TEXT_BRIGHT_BLACK)
	}

	buffer.WriteString(lineNumber)
	buffer.WriteString(constants.FOREGROUND_RESET)
}

func DrawWelcomeMessage(buffer *bytes.Buffer, screenCols int) {
	welcome := "Go editor -- version 0.1"
	welcomelen := len(welcome)
	if welcomelen > screenCols {
		welcomelen = screenCols
	}
	padding := (screenCols - welcomelen) / 2
	if padding > 0 {
		buffer.WriteByte(byte(constants.TILDE))
		padding--
	}
	for padding > 0 {
		buffer.WriteByte(byte(constants.SPACE_RUNE))
		padding--
	}
	buffer.WriteString(welcome)
}

func DrawFileBrowserHeader(buffer *bytes.Buffer, e *config.Editor) {
	// Instructions

	instructionsLines := e.InstructionsLines()
	// Initialize line counter

	for i, line := range instructionsLines {
		cursorPosition := SetCursorPos(i+1, 0)
		buffer.WriteString(cursorPosition)
		buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
		buffer.WriteString(line)
		buffer.WriteString("\n")
	}
}

func DrawFileBrowser(buffer *bytes.Buffer, e *config.Editor, startRow, endRow int) {
	HideCursorIf(buffer, e.FileBrowserActionState.Modifying)
	DrawFileBrowserHeader(buffer, e)
	if e.Cy < len(e.InstructionsLines()) {
		e.Cy = len(e.InstructionsLines())
	}
	var textColor string = "\x1b[32m" // Set text color to green
	var resetColor string = "\x1b[0m" // Reset all terminal attributes to default

	if endRow >= e.ScreenRows-len(e.InstructionsLines()) {
		endRow = e.ScreenRows - len(e.InstructionsLines())
	}

	for i := startRow; i <= endRow; i++ {
		fileRow := i + e.RowOff
		cursorPosition := SetCursorPos(fileRow+1-e.RowOff+len(e.InstructionsLines()), 0)
		buffer.WriteString(cursorPosition)

		if fileRow >= len(e.FileBrowserItems) {
			buffer.WriteString(" ")
		} else {
			item := e.FileBrowserItems[fileRow]
			if item.Type == "directory" {
				buffer.WriteString(textColor + item.Name + "/" + resetColor)
			} else {
				buffer.WriteString(item.Name)
			}
		}

		buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
		buffer.WriteString(constants.ESCAPE_NEW_LINE)
	}
}

func EditorDrawRows(buffer *bytes.Buffer, e *config.Editor, startRow, endRow int) {
	screenCols := e.ScreenCols
	HideCursorIfSearching(buffer, e)

	for i := startRow; i <= endRow; i++ {
		fileRow := i + e.RowOff

		cursorPosition := SetCursorPos(fileRow+1-e.RowOff, 0) // +1 because terminal rows start from 1
		buffer.WriteString(cursorPosition)
		DrawLineNumbers(buffer, fileRow, e)
		cursorPosition = SetCursorPos(fileRow+1-e.RowOff, 6) // +1 because terminal rows start from 1

		if fileRow >= e.CurrentBuffer.NumRows {
			WriteWelcomeIfNoFile(buffer, screenCols, endRow-startRow+1, i, e)
		} else {
			rowLength := e.CurrentBuffer.Rows[fileRow].Length - e.ColOff
			availableScreenCols := screenCols - e.LineNumberWidth
			if fileRow < 5 {
				config.LogToFile(fmt.Sprintf("FileRow: %d, selectedcy: %d, %d", fileRow, e.CurrentBuffer.SelectedCyStart, e.CurrentBuffer.SelectedCyEnd))
			}
			isSelectedRow := fileRow >= e.CurrentBuffer.SelectedCyStart && fileRow <= e.CurrentBuffer.SelectedCyEnd

			if rowLength < 0 {
				rowLength = 0
			}
			if rowLength > availableScreenCols {
				rowLength = availableScreenCols
			}
			if e.ColOff < e.CurrentBuffer.Rows[fileRow].Length {
				highlights := e.CurrentBuffer.Rows[fileRow].Highlighting
				cColor := -1
				for j := 0; j < rowLength; j++ {

					c := e.CurrentBuffer.Rows[fileRow].Chars[e.ColOff+j]
					isSelectedChar := isSelectedRow && j+e.ColOff >= e.CurrentBuffer.SelectedCxStart-e.LineNumberWidth && j+e.ColOff <= e.CurrentBuffer.SelectedCxEnd-e.LineNumberWidth

					hl := highlights[e.ColOff+j]
					if c == ' ' {
						spaceCount := CountSpaces(e, rowLength, j, fileRow)
						if j > constants.TAB_STOP && spaceCount == constants.TAB_STOP {
							AppendTabOrRowIndentBar(e, &j, buffer, fileRow, rowLength)
							j += constants.TAB_STOP - 1
							continue
						}
					}

					if unicode.IsControl(rune(c)) {
						ControlCHandler(buffer, rune(c), cColor)
					} else if isSelectedChar && (fileRow < e.CurrentBuffer.SelectedCyEnd || j+e.ColOff <= e.Cx-e.LineNumberWidth) {
						FormatSelectedTextHandler(buffer, c, &cColor, hl)
					} else if hl == constants.HL_MATCH {
						FormatFindResultHandler(buffer, c)
					} else if hl == constants.HL_NORMAL {
						NormalFormatHandler(buffer, c, cColor)
					} else {
						ColorFormatHandler(buffer, c, &cColor, hl)
					}
				}
				buffer.WriteString(constants.FOREGROUND_RESET)
				cColor = -1
			} else {
				buffer.Write([]byte{})
			}

		}

		buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
		buffer.WriteString(constants.ESCAPE_NEW_LINE)
	}
}

func EditorDrawStatusBar(buf *bytes.Buffer, e *config.Editor) {
	// Set background color for the status bar

	var modeBgColor string
	var modeBold string = "\x1b[1m" // Bold text

	var modeTextColor string = "\x1b[38;5;236m" // Dark gray text
	var modeName string

	switch e.EditorMode {
	case constants.EDITOR_MODE_NORMAL, constants.EDITOR_MODE_FILE_BROWSER:
		modeBgColor = "\x1b[48;5;1m" // Subdued Red background
		modeName = " NORMAL "
	case constants.EDITOR_MODE_VISUAL:
		modeBgColor = "\x1b[48;5;4m" // Subdued Blue background
		modeName = " VISUAL "
	case constants.EDITOR_MODE_INSERT:
		modeBgColor = "\x1b[48;5;2m" // Subdued Green background

		modeName = " INSERT "
	default:
		modeBgColor = "\x1b[48;5;236m" // Default dark gray background
		modeName = "UNKNOWN"
	}
	buf.WriteString(fmt.Sprintf("%s%s%s%-7s\x1b[0m", modeBgColor, modeTextColor, modeBold, modeName))

	buf.WriteString("\x1b[48;5;236m") // Dark gray background

	// File Status Section
	currentRow := e.Cy + 1
	if currentRow > e.CurrentBuffer.NumRows+1 {
		currentRow = e.CurrentBuffer.NumRows + 1
	}

	dirty := ""
	if e.CurrentBuffer.Dirty > 0 {
		dirty = "\x1b[31m(modified)\x1b[39m" // Red color for modified
	}

	status := fmt.Sprintf(" \x1b[32m%.20s\x1b[39m - %d lines %s", e.CurrentBuffer.Name, e.CurrentBuffer.NumRows, dirty) // Green color for filename

	// Right-aligned Status
	rStatus := fmt.Sprintf("%s \x1b[34m|\x1b[39m %d/%d", e.CurrentBuffer.BufferSyntax.FileType, e.Cy+1, e.CurrentBuffer.NumRows) // Blue color for separator

	// Calculate the visible length of status and rStatus (ignoring ANSI codes)
	visibleStatusLen := len(status) - 9   // 9 characters are for ANSI codes in 'status'
	visibleRStatusLen := len(rStatus) - 9 // 9 characters are for ANSI codes in 'rStatus'

	// Calculate the number of spaces needed to fill the gap
	spaceCount := e.ScreenCols - (visibleStatusLen + visibleRStatusLen + 7)

	// Write the status bars
	buf.WriteString(status)

	// Fill the gap with spaces
	for i := 0; i < spaceCount; i++ {
		buf.WriteString(" ")
	}

	// Write the right-aligned status
	buf.WriteString(rStatus)

	// Reset terminal attributes and move to the next line
	buf.WriteString(constants.ESCAPE_RESET_ATTRIBUTES)
	buf.WriteString(constants.ESCAPE_NEW_LINE)
}

func EditorConfirmationPrompt(prompt string, e *config.Editor) bool {
	for {
		EditorSetStatusMessage(e, "%s (y/n)", prompt)
		EditorRefreshScreen(e, constants.INITIAL_REFRESH)
		c, err := ReadKey(e.Reader)
		if err != nil {
			log.Fatal(err)
		}

		if c == 'y' || c == 'Y' {
			EditorSetStatusMessage(e, "")
			return true
		} else if c == 'n' || c == 'N' || c == '\x1b' { // Escape key to cancel
			EditorSetStatusMessage(e, "")
			return false
		}
		// Ignore other keys
	}
}

func EditorPrompt(prompt string, cb func([]rune, rune, *config.Editor, bool), e *config.Editor) []rune {
	buf := []rune{}
	for {
		EditorSetStatusMessage(e, "%s", fmt.Sprintf("%s %s", prompt, string(buf)))
		EditorRefreshScreen(e, constants.INITIAL_REFRESH)
		c, err := ReadKey(e.Reader)
		if err != nil {
			log.Fatal(err)
		}

		if c == constants.DEL_KEY || c == utils.CTRL_KEY('h') || c == constants.BACKSPACE {
			if len(buf) != 0 {
				buf = buf[:len(buf)-1]
				if cb != nil {
					cb(buf, c, e, false)
				}
			}
		} else if c == '\x1b' {
			EditorSetStatusMessage(e, "")
			if cb != nil {
				cb(buf, c, e, false)
			}
			return nil
		} else if c == '\r' {
			if len(buf) != 0 {
				EditorSetStatusMessage(e, "")
				if cb != nil {
					cb(buf, c, e, true)
				}
				return buf
			}
		} else if c != utils.CTRL_KEY('c') && c < 128 {
			buf = append(buf, c)
		}

		if cb != nil {
			cb(buf, c, e, false)
		}
	}
}

func EditorDrawMessageBar(buf *bytes.Buffer, e *config.Editor) {
	buf.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END) // Clear the line
	msgLen := len(e.StatusMsg)
	if msgLen > e.ScreenCols {
		msgLen = e.ScreenCols
	}
	if msgLen > 0 && time.Since(e.StatusMsgTime).Seconds() < 5 {
		buf.WriteString(e.StatusMsg[:msgLen]) // Write the message if within 5 seconds
	}
}
