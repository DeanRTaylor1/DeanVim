package actions

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
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

func EditorRefreshScreen(cfg *config.EditorConfig) {
	EditorScroll(cfg)
	var buffer bytes.Buffer

	buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
	buffer.WriteString(constants.ESCAPE_MOVE_TO_HOME_POS)

	if cfg.Cx < 5 {
		cfg.Cx = 5
	}

	EditorDrawRows(&buffer, cfg)
	EditorDrawStatusBar(&buffer, cfg)
	EditorDrawMessageBar(&buffer, cfg)

	cursorPosition := SetCursorPos((cfg.Cy-cfg.RowOff)+1, (cfg.Cx-cfg.ColOff)+1)
	buffer.WriteString(cursorPosition)

	buffer.WriteString(constants.ESCAPE_SHOW_CURSOR)

	os.Stdout.Write(buffer.Bytes())
}

func EditorDrawRows(buffer *bytes.Buffer, cfg *config.EditorConfig) {
	screenRows := cfg.ScreenRows
	screenCols := cfg.ScreenCols
	if cfg.CurrentBuffer.SearchState.Searching == true {
		buffer.WriteString(constants.ESCAPE_HIDE_CURSOR)
	}
	for i := 0; i < screenRows; i++ {
		fileRow := i + cfg.RowOff
		DrawLineNumbers(buffer, fileRow, cfg)

		if fileRow >= cfg.CurrentBuffer.NumRows {
			if cfg.CurrentBuffer.NumRows == 0 && i == screenRows/3 {
				DrawWelcomeMessage(buffer, screenCols)
			} else {
				buffer.WriteByte(byte(constants.TILDE))
			}
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
						spaceCount := 0
						for k := j; k < j+constants.TAB_STOP; k++ {
							if k >= rowLength || cfg.CurrentBuffer.Rows[fileRow].Chars[cfg.ColOff+k] != ' ' {
								break
							}
							spaceCount++
						}

						if j > constants.TAB_STOP && spaceCount == constants.TAB_STOP {
							nextCharIndex := j + constants.TAB_STOP
							if nextCharIndex < rowLength && cfg.CurrentBuffer.Rows[fileRow].Chars[cfg.ColOff+nextCharIndex] != '}' {
								buffer.WriteString(strings.Repeat(" ", constants.TAB_STOP-1))
								buffer.WriteString(constants.TEXT_BRIGHT_BLACK)
								buffer.WriteString("│")
								buffer.WriteString(constants.FOREGROUND_RESET)
							} else {
								// If the next character is a '}', just append the spaces
								buffer.WriteString(strings.Repeat(" ", constants.TAB_STOP))
							}
							j += constants.TAB_STOP - 1 // Skip ahead
							continue
						}
					}

					if unicode.IsControl(rune(c)) {
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
					} else if hl == constants.HL_MATCH {
						buffer.WriteString(constants.ESCAPE_HIDE_CURSOR)
						buffer.WriteString(constants.FOREGROUND_RESET)
						buffer.WriteString(constants.BACKGROUND_YELLOW)
						buffer.WriteByte(c)
						buffer.WriteString(constants.BACKGROUND_RESET)
					} else if hl == constants.HL_NORMAL {
						if cColor != -1 {
							buffer.WriteString(constants.FOREGROUND_RESET)
							cColor = -1
						}
						buffer.WriteByte(c)
					} else {
						color := int(highlighting.EditorSyntaxToColor(hl))
						if color != cColor {
							buffer.WriteString(fmt.Sprintf("\x1b[%dm", color))
							cColor = color
						}
						buffer.WriteByte(c)
					}
				}
				buffer.WriteString(constants.FOREGROUND_RESET)
			} else {
				buffer.Write([]byte{})
			}

		}
		buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)

		buffer.WriteString(constants.ESCAPE_NEW_LINE)
	}
}

func EditorDrawStatusBar(buf *bytes.Buffer, cfg *config.EditorConfig) {
	buf.WriteString("\x1b[7m")

	currentRow := cfg.Cy + 1
	if currentRow > cfg.CurrentBuffer.NumRows {
		currentRow = cfg.CurrentBuffer.NumRows
	}

	dirty := ""
	if cfg.Dirty > 0 {
		dirty = "(modified)"
	}

	status := fmt.Sprintf("%.20s - %d lines %s", cfg.FileName, cfg.CurrentBuffer.NumRows, dirty)
	rStatus := fmt.Sprintf("%s | %d/%d", cfg.CurrentBuffer.BufferSyntax.FileType, cfg.Cy+1, cfg.CurrentBuffer.NumRows)

	rLen := len(rStatus)
	if len(status) > cfg.ScreenCols {
		status = status[:cfg.ScreenCols-rLen]
	}

	buf.WriteString(status)
	for i := len(status); i < cfg.ScreenCols-rLen; i++ {
		buf.WriteString(" ")
	}

	buf.WriteString(rStatus)
	buf.WriteString(constants.ESCAPE_RESET_ATTRIBUTES)
	buf.WriteString(constants.ESCAPE_NEW_LINE)
}

func EditorPrompt(prompt string, cb func([]rune, rune, *config.EditorConfig), cfg *config.EditorConfig) []rune {
	buf := []rune{}
	for {
		EditorSetStatusMessage(cfg, "%s", fmt.Sprintf("%s %s", prompt, string(buf)))
		EditorRefreshScreen(cfg)
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
