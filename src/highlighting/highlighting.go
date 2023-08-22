package highlighting

import (
	"fmt"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
)

func EditorUpdateSyntax(row *config.Row) {
	hl := make([]byte, row.Length)
	for i, c := range row.Chars {
		hl[i] = constants.HL_NORMAL
		if utils.IsDigit(c) {
			hl[i] = constants.HL_NUMBER
		}
	}
	message := fmt.Sprintf(
		"EditorUpdateSyntax Debug Info:\n"+
			"Length of row.Chars: %d\n"+
			"Length of row.Highlighting: %d\n"+
			"row.Length: %d\n",
		len(row.Chars),
		len(row.Highlighting),
		row.Length,
	)
	config.LogToFile(message)
	row.Highlighting = hl
}

func EditorSyntaxToColor(highlight byte) byte {
	switch highlight {
	case constants.HL_NUMBER:
		return 31
	default:
		return 37
	}
}
