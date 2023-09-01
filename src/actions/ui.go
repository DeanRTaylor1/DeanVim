package actions

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

func DrawLineNumbers(buffer *bytes.Buffer, fileRow int, cfg *config.EditorConfig) {
	relativeLineNumber := int(math.Abs(float64(cfg.Cy - fileRow)))
	lineNumber := fmt.Sprintf("%4d ", relativeLineNumber)

	if fileRow == cfg.Cy {
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

func DrawFileBrowserHeader(buffer *bytes.Buffer) int {
	// Instructions
	instructionsLines := []string{
		"==========================================================",
		"sorted by: date",
		"sort sequence: newest to oldest",
		"Help:<Up> to go up, <Down> to go down, <Enter> to select",
		"==========================================================",
	}

	// Initialize line counter
	lineCount := 0

	for _, line := range instructionsLines {
		cursorPosition := SetCursorPos(lineCount+1, 0)
		buffer.WriteString(cursorPosition)
		buffer.WriteString(line)
		buffer.WriteString("\n")
		lineCount++
	}

	return lineCount
}

func DrawFileBrowser(buffer *bytes.Buffer, cfg *config.EditorConfig, startRow, endRow int) {
	headerLines := DrawFileBrowserHeader(buffer)
	if cfg.Cy < headerLines {
		cfg.Cy = headerLines
	}
	var textColor string = "\x1b[32m" // Set text color to green
	var resetColor string = "\x1b[0m" // Reset all terminal attributes to default

	if endRow >= cfg.ScreenRows-headerLines {
		endRow = cfg.ScreenRows - headerLines
	}

	for i := startRow; i <= endRow; i++ {
		fileRow := i + cfg.RowOff
		cursorPosition := SetCursorPos(fileRow+1-cfg.RowOff+headerLines, 0)
		buffer.WriteString(cursorPosition)

		if fileRow >= len(cfg.FileBrowserItems) {
			buffer.WriteString(" ")
		} else {
			item := cfg.FileBrowserItems[fileRow]
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

func EditorDrawRows(buffer *bytes.Buffer, cfg *config.EditorConfig, startRow, endRow int) {
	screenCols := cfg.ScreenCols
	HideCursorIfSearching(buffer, cfg)

	for i := startRow; i <= endRow; i++ {
		fileRow := i + cfg.RowOff

		cursorPosition := SetCursorPos(fileRow+1-cfg.RowOff, 0) // +1 because terminal rows start from 1
		buffer.WriteString(cursorPosition)
		DrawLineNumbers(buffer, fileRow, cfg)
		cursorPosition = SetCursorPos(fileRow+1-cfg.RowOff, 6) // +1 because terminal rows start from 1

		if fileRow >= cfg.CurrentBuffer.NumRows {
			WriteWelcomeIfNoFile(buffer, screenCols, endRow-startRow+1, i, cfg)
		} else {
			rowLength := cfg.CurrentBuffer.Rows[fileRow].Length - cfg.ColOff
			availableScreenCols := screenCols - cfg.LineNumberWidth
			if rowLength < 0 {
				rowLength = 0
			}
			if rowLength > availableScreenCols {
				rowLength = availableScreenCols
			}
			if cfg.ColOff < cfg.CurrentBuffer.Rows[fileRow].Length {
				highlights := cfg.CurrentBuffer.Rows[fileRow].Highlighting
				cColor := -1
				for j := 0; j < rowLength; j++ {

					c := cfg.CurrentBuffer.Rows[fileRow].Chars[cfg.ColOff+j]

					hl := highlights[cfg.ColOff+j]
					if c == ' ' {
						spaceCount := CountSpaces(cfg, rowLength, j, fileRow)
						if j > constants.TAB_STOP && spaceCount == constants.TAB_STOP {
							AppendTabOrRowIndentBar(cfg, &j, buffer, fileRow, rowLength)
							j += constants.TAB_STOP - 1
							continue
						}
					}
					if unicode.IsControl(rune(c)) {
						ControlCHandler(buffer, rune(c), cColor)
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

func EditorDrawStatusBar(buf *bytes.Buffer, cfg *config.EditorConfig) {
	// Set background color for the status bar

	var modeBgColor string
	var modeBold string = "\x1b[1m" // Bold text

	var modeTextColor string = "\x1b[38;5;236m" // Dark gray text
	var modeName string

	switch cfg.EditorMode {
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
	currentRow := cfg.Cy + 1
	if currentRow > cfg.CurrentBuffer.NumRows+1 {
		currentRow = cfg.CurrentBuffer.NumRows + 1
	}

	dirty := ""
	if cfg.Dirty > 0 {
		dirty = "\x1b[31m(modified)\x1b[39m" // Red color for modified
	}

	status := fmt.Sprintf(" \x1b[32m%.20s\x1b[39m - %d lines %s", cfg.CurrentBuffer.Name, cfg.CurrentBuffer.NumRows, dirty) // Green color for filename

	// Right-aligned Status
	rStatus := fmt.Sprintf("%s \x1b[34m|\x1b[39m %d/%d", cfg.CurrentBuffer.BufferSyntax.FileType, cfg.Cy+1, cfg.CurrentBuffer.NumRows) // Blue color for separator

	// Calculate the visible length of status and rStatus (ignoring ANSI codes)
	visibleStatusLen := len(status) - 9   // 9 characters are for ANSI codes in 'status'
	visibleRStatusLen := len(rStatus) - 9 // 9 characters are for ANSI codes in 'rStatus'

	// Calculate the number of spaces needed to fill the gap
	spaceCount := cfg.ScreenCols - (visibleStatusLen + visibleRStatusLen + 7)

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

func EditorPrompt(prompt string, cb func([]rune, rune, *config.EditorConfig), cfg *config.EditorConfig) []rune {
	buf := []rune{}
	for {
		EditorSetStatusMessage(cfg, "%s", fmt.Sprintf("%s %s", prompt, string(buf)))
		EditorRefreshScreen(cfg, constants.INITIAL_REFRESH)
		c, err := ReadKey(cfg.Reader)
		if err != nil {
			log.Fatal(err)
		}

		if c == constants.DEL_KEY || c == utils.CTRL_KEY('h') || c == constants.BACKSPACE {
			if len(buf) != 0 {
				buf = buf[:len(buf)-1]
				if cb != nil {
					cb(buf, c, cfg)
				}
			}
		} else if c == '\x1b' {
			EditorSetStatusMessage(cfg, "")
			if cb != nil {
				cb(buf, c, cfg)
			}
			return nil
		} else if c == '\r' {
			if len(buf) != 0 {
				EditorSetStatusMessage(cfg, "")
				if cb != nil {
					cb(buf, c, cfg)
				}
				return buf
			}
		} else if c != utils.CTRL_KEY('c') && c < 128 {
			buf = append(buf, c)
		}

		if cb != nil {
			cb(buf, c, cfg)
		}
	}
}

func EditorDrawMessageBar(buf *bytes.Buffer, cfg *config.EditorConfig) {
	buf.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END) // Clear the line
	msgLen := len(cfg.StatusMsg)
	if msgLen > cfg.ScreenCols {
		msgLen = cfg.ScreenCols
	}
	if msgLen > 0 && time.Since(cfg.StatusMsgTime).Seconds() < 5 {
		buf.WriteString(cfg.StatusMsg[:msgLen]) // Write the message if within 5 seconds
	}
}
