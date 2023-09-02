package actions

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func EditorCreateFileCallback(buf []rune, c rune, cfg *config.EditorConfig, trigger bool) {
	buf = append(buf, c)
	if trigger {
		fileName := string(buf)

		// Log the create operation
		config.LogToFile(fmt.Sprintf("FileName to Create: %s", fileName))

		// Show a confirmation prompt
		confirmationMessage := fmt.Sprintf("Are you sure you want to create %s?", fileName)
		if EditorConfirmationPrompt(confirmationMessage, cfg) {
			// Perform the create operation here
			err := EditorCreateFile(cfg, fileName)
			if err != nil {
				EditorSetStatusMessage(cfg, fmt.Sprintf("Failed to create file: %s", err.Error()))
			}
		} else {
			// User said 'no', cancel the operation
			EditorSetStatusMessage(cfg, "Create operation cancelled.")
		}
	}
}

func EditorDeleteFileCallback(buf []rune, c rune, cfg *config.EditorConfig, trigger bool) {
	buf = append(buf, c)
	if trigger {
	}
}

func EditorRenameCallback(buf []rune, c rune, cfg *config.EditorConfig, trigger bool) {
	buf = append(buf, c)
	if trigger {

		oldName := cfg.FileBrowserItems[cfg.Cy-len(cfg.InstructionsLines())].Name
		newName := string(buf)

		// Log the rename operation
		config.LogToFile(fmt.Sprintf("OldName: %s, NewName: %s", oldName, newName))

		// Show a confirmation prompt
		confirmationMessage := "Are you sure you want to rename?"
		if filepath.Ext(oldName) != filepath.Ext(newName) {
			confirmationMessage = "You are changing the file extension. Are you sure you want to rename?"
		}

		if EditorConfirmationPrompt(confirmationMessage, cfg) {
			// Perform the rename operation here
			EditorRenameFile(cfg, oldName, newName)
		} else {
			// User said 'no', cancel the operation
			EditorSetStatusMessage(cfg, "Rename operation cancelled.")
		}
	}
}

func EditorFindCallback(buf []rune, c rune, cfg *config.EditorConfig, trigger bool) {
	if len(cfg.CurrentBuffer.SearchState.SavedHl) > 0 {
		sl := cfg.CurrentBuffer.SearchState.SavedHlLine
		cfg.CurrentBuffer.Rows[sl].Highlighting = cfg.CurrentBuffer.SearchState.SavedHl
	}

	if c == '\r' || c == '\x1b' {
		cfg.CurrentBuffer.SearchState.LastMatch = -1
		cfg.CurrentBuffer.SearchState.Direction = 1
		return
	} else if c == constants.ARROW_RIGHT || c == constants.ARROW_DOWN {
		cfg.CurrentBuffer.SearchState.Direction = 1
	} else if c == constants.ARROW_LEFT || c == constants.ARROW_UP {
		cfg.CurrentBuffer.SearchState.Direction = -1
	} else {
		cfg.CurrentBuffer.SearchState.LastMatch = -1
		cfg.CurrentBuffer.SearchState.Direction = 1
	}

	if cfg.CurrentBuffer.SearchState.LastMatch == -1 {
		cfg.CurrentBuffer.SearchState.Direction = 1
	}
	current := cfg.CurrentBuffer.SearchState.LastMatch
	for i := 0; i < cfg.CurrentBuffer.NumRows; i++ {
		current += cfg.CurrentBuffer.SearchState.Direction
		if current == -1 {
			current = cfg.CurrentBuffer.NumRows - 1
		} else if current == cfg.CurrentBuffer.NumRows {
			current = 0
		}

		row := cfg.CurrentBuffer.Rows[current].Chars
		matchIndex := strings.Index(string(row), string(buf))
		if matchIndex != -1 {
			cfg.CurrentBuffer.SearchState.LastMatch = current
			cfg.Cy = current
			cfg.Cx = matchIndex + cfg.LineNumberWidth
			cfg.CurrentBuffer.SliceIndex = matchIndex
			cfg.RowOff = cfg.CurrentBuffer.NumRows

			cfg.CurrentBuffer.SearchState.SavedHlLine = current
			cfg.CurrentBuffer.SearchState.SavedHl = make([]byte, len(cfg.CurrentBuffer.Rows[cfg.Cy].Highlighting))
			copy(cfg.CurrentBuffer.SearchState.SavedHl, cfg.CurrentBuffer.Rows[cfg.Cy].Highlighting)

			for i := 0; i < len(buf); i++ {
				cfg.CurrentBuffer.Rows[cfg.Cy].Highlighting[matchIndex+i] = constants.HL_MATCH
			}
			break
		}
	}
}

func EditorDelete(cfg *config.EditorConfig) {
	cfg.FileBrowserActionState.Modifying = true
	cfg.FileBrowserActionState.ItemToModify = cfg.FileBrowserItems[cfg.Cy-len(cfg.InstructionsLines())]
	cx := cfg.Cx
	cy := cfg.Cy
	sliceIndex := cfg.CurrentBuffer.SliceIndex
	rowOff := cfg.RowOff
	colOff := cfg.ColOff

	fileName := cfg.FileBrowserItems[cfg.Cy-len(cfg.InstructionsLines())].Name

	// Log the delete operation
	config.LogToFile(fmt.Sprintf("FileName to Delete: %s", fileName))

	// Show a confirmation prompt
	confirmationMessage := fmt.Sprintf("Are you sure you want to delete %s?", fileName)
	if EditorConfirmationPrompt(confirmationMessage, cfg) {
		// Perform the delete operation here
		err := EditorDeleteFile(cfg, fileName)
		if err != nil {
			EditorSetStatusMessage(cfg, fmt.Sprintf("Failed to delete file: %s", err.Error()))
		}
	} else {
		// User said 'no', cancel the operation
		EditorSetStatusMessage(cfg, "Delete operation cancelled.")
	}

	cfg.Cx = cx
	cfg.CurrentBuffer.SliceIndex = sliceIndex
	cfg.Cy = cy
	cfg.RowOff = rowOff
	cfg.ColOff = colOff
	cfg.FileBrowserActionState.Modifying = false
}

func EditorCreate(cfg *config.EditorConfig) {
	cfg.FileBrowserActionState.Modifying = true
	cx := cfg.Cx
	cy := cfg.Cy
	sliceIndex := cfg.CurrentBuffer.SliceIndex
	rowOff := cfg.RowOff
	colOff := cfg.ColOff

	EditorPrompt("Create File: ", EditorCreateFileCallback, cfg)

	cfg.Cx = cx
	cfg.CurrentBuffer.SliceIndex = sliceIndex
	cfg.Cy = cy
	cfg.RowOff = rowOff
	cfg.ColOff = colOff
	cfg.FileBrowserActionState.Modifying = false
}

func EditorRename(cfg *config.EditorConfig) {
	cfg.FileBrowserActionState.Modifying = true
	cfg.FileBrowserActionState.ItemToModify = cfg.FileBrowserItems[cfg.Cy-len(cfg.InstructionsLines())]
	cx := cfg.Cx
	cy := cfg.Cy
	sliceIndex := cfg.CurrentBuffer.SliceIndex
	rowOff := cfg.RowOff
	colOff := cfg.ColOff

	EditorPrompt(fmt.Sprintf("Rename File: %s", cfg.FileBrowserActionState.ItemToModify.Name), EditorRenameCallback, cfg)

	cfg.Cx = cx
	cfg.CurrentBuffer.SliceIndex = sliceIndex
	cfg.Cy = cy
	cfg.RowOff = rowOff
	cfg.ColOff = colOff
	cfg.FileBrowserActionState.Modifying = false
}

func EditorFind(cfg *config.EditorConfig) {
	cfg.CurrentBuffer.SearchState.Searching = true
	cx := cfg.Cx
	sliceIndex := cfg.CurrentBuffer.SliceIndex
	cy := cfg.Cy
	rowOff := cfg.RowOff
	colOff := cfg.ColOff

	query := EditorPrompt("Search: (ESC to cancel)", EditorFindCallback, cfg)

	if query == nil {
		cfg.Cx = cx
		cfg.CurrentBuffer.SliceIndex = sliceIndex
		cfg.Cy = cy
		cfg.RowOff = rowOff
		cfg.ColOff = colOff
	}
	cfg.CurrentBuffer.SearchState.Searching = false
}
