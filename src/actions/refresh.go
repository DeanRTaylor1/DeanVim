package actions

import (
	"bytes"
	"fmt"
	"os"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func FullRefresh(cfg *config.EditorConfig, buffer *bytes.Buffer) {
	buffer.WriteString(constants.ESCAPE_MOVE_TO_HOME_POS)
	buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
	if cfg.EditorMode == constants.EDITOR_MODE_FILE_BROWSER {
		DrawFileBrowser(buffer, cfg, 0, cfg.ScreenRows)
	} else {
		EditorDrawRows(buffer, cfg, 0, cfg.ScreenRows)
	}
}

func DrawAllLineNumbers(buffer *bytes.Buffer, cfg *config.EditorConfig) {
	for i := 1; i <= cfg.ScreenRows; i++ {
		fileRow := i + cfg.RowOff - 1
		if fileRow < 0 {
			continue
		}
		cursorPosition := SetCursorPos(i+1, 1)
		buffer.WriteString(cursorPosition)
		DrawLineNumbers(buffer, fileRow, cfg)
	}
}

func PartialRefresh(cfg *config.EditorConfig, buffer *bytes.Buffer, startRow, endRow int) {
	if cfg.EditorMode != constants.EDITOR_MODE_FILE_BROWSER {
		DrawAllLineNumbers(buffer, cfg)
	}
	cursorPosition := SetCursorPos(startRow+1, 6)
	if cfg.EditorMode == constants.EDITOR_MODE_FILE_BROWSER {
		cursorPosition = SetCursorPos(startRow+6, 0)
	}
	buffer.WriteString(cursorPosition)
	buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)

	if cfg.EditorMode == constants.EDITOR_MODE_FILE_BROWSER {
		DrawFileBrowser(buffer, cfg, startRow, endRow)
	} else {
		EditorDrawRows(buffer, cfg, startRow, endRow)
	}
}

func SingleLineRefresh(cfg *config.EditorConfig, buffer *bytes.Buffer, startRow, endRow int) {
	cursorPosition := SetCursorPos(startRow+1, 0)
	buffer.WriteString(cursorPosition)
	buffer.WriteString(constants.ESCAPE_CLEAR_TO_LINE_END)
	EditorDrawRows(buffer, cfg, startRow, endRow)
}

func EditorRefreshScreen(cfg *config.EditorConfig, lastKeyPress rune) {
	var buffer bytes.Buffer
	EditorScroll(cfg)
	buffer.WriteString(constants.ESCAPE_HIDE_CURSOR)

	// Cursor position adjustment logic for non file browser modes
	if cfg.EditorMode != constants.EDITOR_MODE_FILE_BROWSER && cfg.Cx < 5 {
		cfg.Cx = 5
	}

	if cfg.EditorMode == constants.EDITOR_MODE_FILE_BROWSER && cfg.Cy < 5 {
		cfg.Cy = 5
	}

	if cfg.EditorMode != constants.EDITOR_MODE_FILE_BROWSER && (cfg.Cx >= cfg.ScreenCols-cfg.LineNumberWidth || cfg.Cy >= cfg.ScreenRows || cfg.Cx-cfg.LineNumberWidth < cfg.ColOff || cfg.Cy-cfg.RowOff < 0 || cfg.Cx-cfg.ColOff == 5) {
		FullRefresh(cfg, &buffer)
	} else {
		startRow := cfg.Cy - 2
		endRow := cfg.Cy + 2
		if startRow < 0 {
			startRow = 0
		}
		switch lastKeyPress {
		case constants.INITIAL_REFRESH, constants.ENTER_KEY, constants.BACKSPACE, constants.DEL_KEY, utils.CTRL_KEY(lastKeyPress), constants.PAGE_DOWN, constants.PAGE_UP:
			FullRefresh(cfg, &buffer)
		case constants.ARROW_DOWN, constants.ARROW_UP, constants.ARROW_LEFT, constants.ARROW_RIGHT:
			config.LogToFile(fmt.Sprintf("startRow: %d", startRow))
			PartialRefresh(cfg, &buffer, startRow, endRow)
		default:
			SingleLineRefresh(cfg, &buffer, 0, cfg.Cy)
		}
	}

	// Draw status and message bars
	statusBarPosition := SetCursorPos(cfg.ScreenRows+1, 0)
	buffer.WriteString(statusBarPosition)
	EditorDrawStatusBar(&buffer, cfg)
	EditorDrawMessageBar(&buffer, cfg)

	// Set cursor position
	cursorPosition := SetCursorPos((cfg.Cy-cfg.RowOff)+1, (cfg.Cx-cfg.ColOff)+1)
	buffer.WriteString(cursorPosition)

	buffer.WriteString(constants.ESCAPE_SHOW_CURSOR)

	// Write to stdout
	os.Stdout.Write(buffer.Bytes())
}
