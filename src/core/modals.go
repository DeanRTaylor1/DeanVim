package core

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/fuzzy"
)

func DrawTopLabel(buffer *bytes.Buffer, startX, startY, width int, label, bgColor, textColor string) {
	labelStart := startX + (width-len(label))/2
	buffer.WriteString(SetCursorPos(startY, labelStart))
	buffer.WriteString(bgColor)   // Highlight background
	buffer.WriteString(textColor) // Text color
	buffer.WriteString(label)
	buffer.WriteString(constants.BACKGROUND_RESET)        // Reset background
	buffer.WriteString(constants.FOREGROUND_RESET)        // Reset foreground
	buffer.WriteString(constants.ESCAPE_RESET_ATTRIBUTES) // Reset all attributes
}

func DrawContent(buffer *bytes.Buffer, startX, startY, width, height int, e *config.Editor) {
	results, ok := e.Modal.Results.([]string)
	if !ok {
		return
	}

	for i := 0; i < 3 && i < len(results); i++ {
		config.LogToFile(fmt.Sprintf("Result %d: %s", i+1, results[i]))
	}

	for i := 1; i < height-5; i++ {
		cursorPos := SetCursorPos(startY+i, startX)
		config.LogToFile(fmt.Sprintf("Cursor Position Before Vertical Line: %s", cursorPos))

		buffer.WriteString(cursorPos)
		buffer.WriteString(constants.VERTICAL_LINE)

		dataIndex := i - 1 + e.Modal.DataRowOffset

		if dataIndex == e.Modal.ItemIndex {
			buffer.WriteString(constants.BACKGROUND_BRIGHT_BLACK)
			buffer.WriteString(constants.TEXT_BOLD)
		}

		if dataIndex < len(results) {
			str := results[dataIndex]
			strLen := len(str)
			maxStrLen := width - 2
			config.LogToFile(fmt.Sprintf("String Length: %d, Max String Length: %d, Width: %d", strLen, maxStrLen, width))

			if strLen > maxStrLen {
				str = str[:maxStrLen]
			}

			buffer.WriteString(str)

			remainingSpace := maxStrLen - len(str)
			config.LogToFile(fmt.Sprintf("Remaining Space: %d", remainingSpace))

			buffer.WriteString(strings.Repeat(" ", remainingSpace))
		} else {
			buffer.WriteString(strings.Repeat(" ", width-2))
		}

		if dataIndex == e.Modal.ItemIndex {
			buffer.WriteString(constants.BACKGROUND_RESET)
			buffer.WriteString(constants.FOREGROUND_RESET)
			buffer.WriteString(constants.ESCAPE_RESET_ATTRIBUTES)
		}

		cursorPosBeforeLastVerticalLine := SetCursorPos(startY+i, startX+width-1)
		config.LogToFile(fmt.Sprintf("Cursor Position Just Before Last Vertical Line: %s", cursorPosBeforeLastVerticalLine))
		buffer.WriteString(constants.VERTICAL_LINE)
	}
}

func DrawFuzzyContent(buffer *bytes.Buffer, startX, startY, width, height int, e *config.Editor) {
	results, ok := e.Modal.Results.(fuzzy.Matches)
	if !ok {
		return
	}
	for i := 1; i < height-5; i++ {
		buffer.WriteString(SetCursorPos(startY+i, startX))
		buffer.WriteString(constants.VERTICAL_LINE) // Vertical line

		// Check if the index exists in e.Modal.Results
		dataIndex := i - 1 + e.Modal.DataRowOffset // Adjusting the index

		if dataIndex == e.Modal.ItemIndex {
			// Highlight the entire line
			buffer.WriteString(constants.BACKGROUND_BRIGHT_BLACK)
			buffer.WriteString(constants.TEXT_BOLD)
		}

		if dataIndex < len(results) {
			// Write the data at the index
			str := results[dataIndex].Str
			matchedIndexes := results[dataIndex].MatchedIndexes

			for j, char := range str {
				if contains(matchedIndexes, j) {
					buffer.WriteString(constants.TEXT_BLUE)
				}
				buffer.WriteString(string(char))
				if contains(matchedIndexes, j) {
					buffer.WriteString(constants.FOREGROUND_RESET)
				}
			}

			// Fill the remaining space with empty characters
			remainingSpace := width - 2 - len(str)
			buffer.WriteString(strings.Repeat(" ", remainingSpace))
		} else {
			// If the index doesn't exist, fill the entire space with empty characters
			buffer.WriteString(strings.Repeat(" ", width-2))
		}

		if dataIndex == e.Modal.ItemIndex {
			// End the highlight
			buffer.WriteString(constants.BACKGROUND_RESET)
			buffer.WriteString(constants.FOREGROUND_RESET)
			buffer.WriteString(constants.ESCAPE_RESET_ATTRIBUTES)
		}

		buffer.WriteString(constants.VERTICAL_LINE) // Vertical line
	}
}

func DrawContentArea(buffer *bytes.Buffer, startX, startY, width, height int, e *config.Editor) {
	// Draw the top border with rounded corners
	if !e.Modal.ModalDrawn {
		// Draw the top border with rounded corners
		buffer.WriteString(SetCursorPos(startY, startX))
		buffer.WriteString(constants.LEFT_TOP_CORNER)                          // Left-top corner
		buffer.WriteString(strings.Repeat(constants.HORIZONTAL_LINE, width-2)) // Horizontal line
		buffer.WriteString(constants.RIGHT_TOP_CORNER)                         // Right-top corner
	}

	switch e.Modal.Type {
	case config.MODAL_TYPE_FUZZY:
		DrawFuzzyContent(buffer, startX, startY, width, height, e)

	default:
		DrawContent(buffer, startX, startY, width, height, e)
	}

	if !e.Modal.ModalDrawn {
		// Draw the bottom border with rounded corners
		buffer.WriteString(SetCursorPos(startY+height-5, startX))
		buffer.WriteString(constants.LEFT_BOTTOM_CORNER)                       // Left-bottom corner
		buffer.WriteString(strings.Repeat(constants.HORIZONTAL_LINE, width-2)) // Horizontal line
		buffer.WriteString(constants.RIGHT_BOTTOM_CORNER)                      // Right-bottom corner
	}

	e.Modal.ModalDrawn = true
}

// Helper function to check if a slice contains an element
func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func DrawSearchBox(buffer *bytes.Buffer, startX, startY, width int, e *config.Editor) {
	// Draw the top border of the search box with rounded corners

	buffer.WriteString(SetCursorPos(startY, startX))
	buffer.WriteString(constants.LEFT_TOP_CORNER)                          // Left-top corner of the search box
	buffer.WriteString(strings.Repeat(constants.HORIZONTAL_LINE, width-2)) // Horizontal line
	buffer.WriteString(constants.RIGHT_TOP_CORNER)                         // Right-top corner of the search box

	// Draw the sides of the search box and include the ModalInput text
	buffer.WriteString(SetCursorPos(startY+1, startX))
	buffer.WriteString(constants.VERTICAL_LINE) // Left vertical line

	// Convert ModalInput to string and write it inside the search box
	inputText := string(e.Modal.ModalInput[e.Modal.SearchColOffset:])
	buffer.WriteString(inputText)

	// Fill the remaining space with empty characters
	remainingSpace := width - 2 - len(inputText)
	buffer.WriteString(strings.Repeat(" ", remainingSpace))

	buffer.WriteString(constants.VERTICAL_LINE) // Right vertical line

	// Draw the bottom border of the search box with rounded corners
	buffer.WriteString(SetCursorPos(startY+2, startX))
	buffer.WriteString(constants.LEFT_BOTTOM_CORNER)                       // Left-bottom corner of the search box
	buffer.WriteString(strings.Repeat(constants.HORIZONTAL_LINE, width-2)) // Horizontal line
	buffer.WriteString(constants.RIGHT_BOTTOM_CORNER)                      // Right-bottom corner of the search box
}

func EditorDrawModal(buffer *bytes.Buffer, e *config.Editor) string {
	// Calculate the dimensions and position of the modal
	modalAvailableWidth := e.ScreenCols * 85 / 100
	modalAvailableHeight := e.ScreenRows * 90 / 100
	modalWidth := modalAvailableWidth
	modalHeight := modalAvailableHeight
	startX := (e.ScreenCols - modalWidth) / 2
	startY := ((e.ScreenRows - modalHeight) / 2) + 2

	label1 := "Results"
	label1Start := startX + (modalWidth-len(label1))/2

	DrawContentArea(buffer, startX, startY, modalWidth, modalHeight, e)
	DrawTopLabel(buffer, label1Start, startY, len(label1), label1, constants.BACKGROUND_BLUE, constants.TEXT_BLACK)

	searchBoxStartY := modalAvailableHeight + 1
	searchBoxWidth := modalWidth

	DrawSearchBox(buffer, startX, searchBoxStartY, searchBoxWidth, e)

	label2 := "Search"
	label2Start := startX + (searchBoxWidth-len(label2))/2

	DrawTopLabel(buffer, label2Start, searchBoxStartY, len(label2), label2, constants.BACKGROUND_YELLOW, constants.TEXT_BLACK)

	cursorX := startX + 1 + e.Modal.CursorPosition
	cursorY := searchBoxStartY + 1
	return SetCursorPos(cursorY, cursorX)
}

func insertCharModalInput(char rune, e *config.Editor) {
	e.Modal.ModalInput = append(e.Modal.ModalInput, 0)
	copy(e.Modal.ModalInput[e.Modal.CursorPosition+1:], e.Modal.ModalInput)
	e.Modal.ModalInput[e.Modal.CursorPosition] = byte(char)

	e.Modal.CursorPosition++
}
