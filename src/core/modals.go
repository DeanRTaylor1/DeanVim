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

	// Draw the top label for the first modal
	label1 := "Results"
	label1Start := startX + (modalWidth-len(label1))/2
	buffer.WriteString(SetCursorPos(startY, label1Start))
	buffer.WriteString(constants.BACKGROUND_BLUE) // Highlight background
	buffer.WriteString(constants.TEXT_BLACK)
	buffer.WriteString(label1)
	buffer.WriteString(constants.BACKGROUND_RESET)        // Reset background
	buffer.WriteString(constants.FOREGROUND_RESET)        // Reset foreground
	buffer.WriteString(constants.ESCAPE_RESET_ATTRIBUTES) // Reset all attributes

	for i := 1; i < modalHeight-5; i++ {
		buffer.WriteString(SetCursorPos(startY+i, startX))
		buffer.WriteString(constants.VERTICAL_LINE) // Vertical line

		// Check if the index exists in e.Modals.Results
		dataIndex := i - 1 + e.Modal.DataRowOffset // Adjusting the index

		if dataIndex == e.Modal.ItemIndex {
			// Highlight the entire line
			buffer.WriteString(constants.BACKGROUND_BRIGHT_BLACK)
			buffer.WriteString(constants.TEXT_BOLD)
			// buffer.WriteString(constants.TEXT_SPECIAL_DARK_GREY)
		}
		if dataIndex < len(e.Modal.Results) {
			// Write the data at the index
			buffer.WriteString(e.Modal.Results[dataIndex])

			// Fill the remaining space with empty characters
			remainingSpace := modalWidth - 2 - len(e.Modal.Results[dataIndex])
			buffer.WriteString(strings.Repeat(" ", remainingSpace))
		} else {
			// If the index doesn't exist, fill the entire space with empty characters
			buffer.WriteString(strings.Repeat(" ", modalWidth-2))
		}

		if dataIndex == e.Modal.ItemIndex {
			// End the highlight
			buffer.WriteString(constants.BACKGROUND_RESET)
			buffer.WriteString(constants.FOREGROUND_RESET)
			buffer.WriteString(constants.ESCAPE_RESET_ATTRIBUTES)
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

	label2 := "Search"

	label2Start := startX + (searchBoxWidth-len(label2))/2
	buffer.WriteString(SetCursorPos(searchBoxStartY, label2Start))
	buffer.WriteString(constants.BACKGROUND_YELLOW) // Highlight background
	buffer.WriteString(constants.TEXT_BLACK)
	buffer.WriteString(label2)
	buffer.WriteString(constants.BACKGROUND_RESET)        // Reset background
	buffer.WriteString(constants.FOREGROUND_RESET)        // Reset foreground
	buffer.WriteString(constants.ESCAPE_RESET_ATTRIBUTES) // Reset all attributes

	// Draw the sides of the search box and include the ModalInput text
	buffer.WriteString(SetCursorPos(searchBoxStartY+1, startX))
	buffer.WriteString(constants.VERTICAL_LINE) // Left vertical line

	// Convert ModalInput to string and write it inside the search box
	inputText := string(e.Modal.ModalInput[e.Modal.SearchColOffset:])
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
