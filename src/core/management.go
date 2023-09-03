package core

import (
	"fmt"
	"path/filepath"

	"github.com/deanrtaylor1/go-editor/config"
)

func EditorCreateFileCallback(buf []rune, c rune, e *config.Editor, trigger bool) {
	buf = append(buf, c)
	if trigger {
		fileName := string(buf)

		// Log the create operation

		// Show a confirmation prompt
		confirmationMessage := fmt.Sprintf("Are you sure you want to create %s?", fileName)
		if EditorConfirmationPrompt(confirmationMessage, e) {
			// Perform the create operation here
			err := EditorCreateFile(e, fileName)
			if err != nil {
				EditorSetStatusMessage(e, fmt.Sprintf("Failed to create file: %s", err.Error()))
			}
		} else {
			// User said 'no', cancel the operation
			EditorSetStatusMessage(e, "Create operation cancelled.")
		}
	}
}

func EditorRenameCallback(buf []rune, c rune, e *config.Editor, trigger bool) {
	buf = append(buf, c)
	if trigger {

		oldName := e.FileBrowserItems[e.Cy-len(e.InstructionsLines())].Name
		newName := string(buf)

		// Log the rename operation

		// Show a confirmation prompt
		confirmationMessage := "Are you sure you want to rename?"
		if filepath.Ext(oldName) != filepath.Ext(newName) {
			confirmationMessage = "You are changing the file extension. Are you sure you want to rename?"
		}

		if EditorConfirmationPrompt(confirmationMessage, e) {
			// Perform the rename operation here
			EditorRenameFile(e, oldName, newName)
		} else {
			// User said 'no', cancel the operation
			EditorSetStatusMessage(e, "Rename operation cancelled.")
		}
	}
}

func EditorDelete(e *config.Editor) {
	e.FileBrowserActionState.Modifying = true
	e.FileBrowserActionState.ItemToModify = e.FileBrowserItems[e.Cy-len(e.InstructionsLines())]
	cx := e.Cx
	cy := e.Cy
	sliceIndex := e.CurrentBuffer.SliceIndex
	rowOff := e.RowOff
	colOff := e.ColOff

	fileName := e.FileBrowserItems[e.Cy-len(e.InstructionsLines())].Name

	// Show a confirmation prompt
	confirmationMessage := fmt.Sprintf("Are you sure you want to delete %s?", fileName)
	if EditorConfirmationPrompt(confirmationMessage, e) {
		// Perform the delete operation here
		err := EditorDeleteFile(e, fileName)
		if err != nil {
			EditorSetStatusMessage(e, fmt.Sprintf("Failed to delete file: %s", err.Error()))
		}
	} else {
		// User said 'no', cancel the operation
		EditorSetStatusMessage(e, "Delete operation cancelled.")
	}

	e.Cx = cx
	e.CurrentBuffer.SliceIndex = sliceIndex
	e.Cy = cy
	e.RowOff = rowOff
	e.ColOff = colOff
	e.FileBrowserActionState.Modifying = false
}

func EditorCreate(e *config.Editor) {
	e.FileBrowserActionState.Modifying = true
	cx := e.Cx
	cy := e.Cy
	sliceIndex := e.CurrentBuffer.SliceIndex
	rowOff := e.RowOff
	colOff := e.ColOff

	EditorPrompt("Create File: ", EditorCreateFileCallback, e)

	e.Cx = cx
	e.CurrentBuffer.SliceIndex = sliceIndex
	e.Cy = cy
	e.RowOff = rowOff
	e.ColOff = colOff
	e.FileBrowserActionState.Modifying = false
}

func EditorRename(e *config.Editor) {
	e.FileBrowserActionState.Modifying = true
	e.FileBrowserActionState.ItemToModify = e.FileBrowserItems[e.Cy-len(e.InstructionsLines())]
	cx := e.Cx
	cy := e.Cy
	sliceIndex := e.CurrentBuffer.SliceIndex
	rowOff := e.RowOff
	colOff := e.ColOff

	EditorPrompt(fmt.Sprintf("Rename File: %s", e.FileBrowserActionState.ItemToModify.Name), EditorRenameCallback, e)

	e.Cx = cx
	e.CurrentBuffer.SliceIndex = sliceIndex
	e.Cy = cy
	e.RowOff = rowOff
	e.ColOff = colOff
	e.FileBrowserActionState.Modifying = false
}
