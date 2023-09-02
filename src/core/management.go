package core

import (
	"fmt"
	"path/filepath"

	"github.com/deanrtaylor1/go-editor/config"
)

func EditorCreateFileCallback(buf []rune, c rune, cfg *config.Editor, trigger bool) {
	buf = append(buf, c)
	if trigger {
		fileName := string(buf)

		// Log the create operation

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

func EditorRenameCallback(buf []rune, c rune, cfg *config.Editor, trigger bool) {
	buf = append(buf, c)
	if trigger {

		oldName := cfg.FileBrowserItems[cfg.Cy-len(cfg.InstructionsLines())].Name
		newName := string(buf)

		// Log the rename operation

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

func EditorDelete(cfg *config.Editor) {
	cfg.FileBrowserActionState.Modifying = true
	cfg.FileBrowserActionState.ItemToModify = cfg.FileBrowserItems[cfg.Cy-len(cfg.InstructionsLines())]
	cx := cfg.Cx
	cy := cfg.Cy
	sliceIndex := cfg.CurrentBuffer.SliceIndex
	rowOff := cfg.RowOff
	colOff := cfg.ColOff

	fileName := cfg.FileBrowserItems[cfg.Cy-len(cfg.InstructionsLines())].Name

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

func EditorCreate(cfg *config.Editor) {
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

func EditorRename(cfg *config.Editor) {
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
