package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/deanrtaylor1/go-editor/constants"
)

const logging = true

func LogToFile(message string) {
	if !logging {
		return
	}

	// Open the log file in append mode, or create it if it doesn't exist
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a log entry with the current time
	logEntry := fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), message)

	// Write the log entry to the file
	if _, err := file.WriteString(logEntry); err != nil {
		log.Fatal(err)
	}
}

type EditorAction struct {
	ActionType   int
	Row          Row
	Index        int
	PrevRow      interface{}
	Cx           int
	RedoFunction func()
}

type FileBrowserActionState struct {
	Modifying    bool
	ItemToModify FileBrowserItem
}

type SearchState struct {
	LastMatch   int
	Direction   int
	SavedHlLine int
	SavedHl     []byte
	Searching   bool
}

type Row struct {
	CharAdjustment   int
	IndentationLevel int
	Idx              int
	Chars            []byte
	Length           int
	Highlighting     []byte
	HlOpenComment    bool
	Tabs             []byte
}

type FileBrowserItem struct {
	Name       string
	Path       string
	Type       string
	Extension  string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func NewBufferSyntax() *BufferSyntax {
	return &BufferSyntax{
		FileType:               "",
		Flags:                  0,
		SingleLineCommentStart: "",
		Syntaxes:               constants.Syntaxes,
	}
}

func NewSearchState() *SearchState {
	return &SearchState{
		LastMatch:   -1,
		Direction:   1,
		SavedHlLine: 0,
		SavedHl:     []byte{},
		Searching:   false,
	}
}

func NewRow() *Row {
	return &Row{
		Idx:              0,
		CharAdjustment:   0,
		IndentationLevel: 0,
		Chars:            []byte{},
		Length:           0,
		Highlighting:     []byte{},
		HlOpenComment:    false,
	}
}

func (r *Row) DeepCopy() *Row {
	// Create a new Row object and copy over the simple fields
	newRow := &Row{
		Idx:              r.Idx,
		CharAdjustment:   r.CharAdjustment,
		IndentationLevel: r.IndentationLevel,
		Length:           r.Length,
		HlOpenComment:    r.HlOpenComment,
	}

	newRow.Chars = make([]byte, len(r.Chars))
	copy(newRow.Chars, r.Chars)

	newRow.Highlighting = make([]byte, len(r.Highlighting))
	copy(newRow.Highlighting, r.Highlighting)

	newRow.Tabs = make([]byte, len(r.Tabs))
	copy(newRow.Tabs, r.Tabs)

	return newRow
}

func GetWindowSize(cfg *Editor) error {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	cfg.ScreenCols = width
	cfg.ScreenRows = height
	return nil
}
