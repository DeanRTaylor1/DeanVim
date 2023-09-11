package core

import (
	"bytes"
	"strings"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func EditorDrawModal(buffer *bytes.Buffer, e *config.Editor) string {
	// Calculate the dimensions and position of the modal
	modalWidth := e.ScreenCols * 85 / 100
	modalHeight := e.ScreenRows * 90 / 100
	startX := (e.ScreenCols - modalWidth) / 2
	startY := ((e.ScreenRows - modalHeight) / 2) + 2

	// Draw the top border with rounded corners
	buffer.WriteString(SetCursorPos(startY, startX))
	buffer.WriteString(constants.LEFT_TOP_CORNER)                               // Left-top corner
	buffer.WriteString(strings.Repeat(constants.HORIZONTAL_LINE, modalWidth-2)) // Horizontal line
	buffer.WriteString(constants.RIGHT_TOP_CORNER)                              // Right-top corner

	for i := 1; i < modalHeight-5; i++ {
		buffer.WriteString(SetCursorPos(startY+i, startX))
		buffer.WriteString(constants.VERTICAL_LINE) // Vertical line

		// Check if the index exists in e.Modal.Data
		dataIndex := i - 1 // Adjusting the index
		if dataIndex < len(e.Modal.Data) {
			// Write the data at the index
			buffer.WriteString(e.Modal.Data[dataIndex])

			// Fill the remaining space with empty characters
			remainingSpace := modalWidth - 2 - len(e.Modal.Data[dataIndex])
			buffer.WriteString(strings.Repeat(" ", remainingSpace))
		} else {
			// If the index doesn't exist, fill the entire space with empty characters
			buffer.WriteString(strings.Repeat(" ", modalWidth-2))
		}

		buffer.WriteString(constants.VERTICAL_LINE) // Vertical line
	}

	// Draw the bottom border with rounded corners
	buffer.WriteString(SetCursorPos(startY+modalHeight-5, startX))
	buffer.WriteString(constants.LEFT_BOTTOM_CORNER)                            // Left-bottom corner
	buffer.WriteString(strings.Repeat(constants.HORIZONTAL_LINE, modalWidth-2)) // Horizontal line
	buffer.WriteString(constants.RIGHT_BOTTOM_CORNER)                           // Right-bottom corner

	searchBoxStartY := modalHeight + 1 // 3 lines from the bottom of the modal
	searchBoxWidth := modalWidth       // 2 spaces padding on each side

	// Draw the top border of the search box with rounded corners
	buffer.WriteString(SetCursorPos(searchBoxStartY, startX))
	buffer.WriteString(constants.LEFT_TOP_CORNER)                                   // Left-top corner of the search box
	buffer.WriteString(strings.Repeat(constants.HORIZONTAL_LINE, searchBoxWidth-2)) // Horizontal line
	buffer.WriteString(constants.RIGHT_TOP_CORNER)                                  // Right-top corner of the search box
	// ... (existing code)

	// Draw the sides of the search box and include the ModalInput text
	buffer.WriteString(SetCursorPos(searchBoxStartY+1, startX))
	buffer.WriteString(constants.VERTICAL_LINE) // Left vertical line

	// Convert ModalInput to string and write it inside the search box
	inputText := string(e.Modal.ModalInput)
	buffer.WriteString(inputText)

	// Fill the remaining space with empty characters
	remainingSpace := searchBoxWidth - 2 - len(inputText)
	buffer.WriteString(strings.Repeat(" ", remainingSpace))

	buffer.WriteString(constants.VERTICAL_LINE) // Right vertical line

	// Draw the bottom border of the search box with rounded corners
	buffer.WriteString(SetCursorPos(searchBoxStartY+2, startX))
	buffer.WriteString(constants.LEFT_BOTTOM_CORNER)                                // Left-bottom corner of the search box
	buffer.WriteString(strings.Repeat(constants.HORIZONTAL_LINE, searchBoxWidth-2)) // Horizontal line
	buffer.WriteString(constants.RIGHT_BOTTOM_CORNER)                               // Right-bottom corner of the search box

	// Place the cursor inside the search box
	cursorX := startX + 1 + e.Modal.CursorPosition
	cursorY := searchBoxStartY + 1 // Inside the search box
	return SetCursorPos(cursorY, cursorX)
}

func insertCharModalInput(char rune, e *config.Editor) {
	e.Modal.ModalInput = append(e.Modal.ModalInput, 0)
	copy(e.Modal.ModalInput[e.Modal.CursorPosition+1:], e.Modal.ModalInput)
	e.Modal.ModalInput[e.Modal.CursorPosition] = byte(char)

	e.Modal.CursorPosition++
}
