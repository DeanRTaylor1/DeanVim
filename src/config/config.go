package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
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
	ActionType int
	Row        Row
	Index      int
	PrevRow    interface{}
	Cx         int
}

type BufferSyntax struct {
	FileType               string
	Flags                  int
	SingleLineCommentStart string
	MultiLineCommentStart  string
	MultiLineCommentEnd    string
	Keywords               []string
	Syntaxes               []constants.SyntaxHighlighting
}

type SearchState struct {
	LastMatch   int
	Direction   int
	SavedHlLine int
	SavedHl     []byte
	Searching   bool
}

type Buffer struct {
	Rows         []Row
	NumRows      int
	SearchState  *SearchState
	BufferSyntax *BufferSyntax
	UndoStack    []EditorAction
	RedoStack    []EditorAction
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

type EditorConfig struct {
	Cx              int
	Cy              int
	SliceIndex      int
	LineNumberWidth int
	ScreenRows      int
	ScreenCols      int
	TerminalState   *term.State
	CurrentBuffer   *Buffer
	RowOff          int
	ColOff          int
	FileName        string
	StatusMsg       string
	StatusMsgTime   time.Time
	Dirty           int
	QuitTimes       int
	Reader          *bufio.Reader
	FirstRead       bool
	UndoHistory     int
}

func (b *Buffer) PopUndo() (EditorAction, bool) {
	if len(b.UndoStack) == 0 {
		return EditorAction{}, false
	}
	lastAction := b.UndoStack[len(b.UndoStack)-1]
	b.UndoStack = b.UndoStack[:len(b.UndoStack)-1]
	return lastAction, true
}

func (b *Buffer) PopRedo() (EditorAction, bool) {
	if len(b.RedoStack) == 0 {
		return EditorAction{}, false
	}
	lastAction := b.RedoStack[len(b.RedoStack)-1]
	b.RedoStack = b.RedoStack[:len(b.RedoStack)-1]
	return lastAction, true
}

func (b *Buffer) AppendRedo(action EditorAction, maxUndoHistory int) {
	if len(b.RedoStack) >= maxUndoHistory {
		b.RedoStack = b.RedoStack[1:]
	}
	b.RedoStack = append(b.RedoStack, action)
}

func (b *Buffer) AppendUndo(action EditorAction, maxUndoHistory int) {
	if len(b.UndoStack) >= maxUndoHistory {
		b.UndoStack = b.UndoStack[1:]
	}
	b.UndoStack = append(b.UndoStack, action)
}

func (c *EditorConfig) GetAdjustedCx() int {
	adjustedCx := utils.Max(0, c.Cx-c.LineNumberWidth)
	adjustedCx = utils.Min(adjustedCx, len(c.CurrentBuffer.Rows[c.Cy].Chars))
	return adjustedCx
}

func (c *EditorConfig) GetCurrentRow() *Row {
	return &c.CurrentBuffer.Rows[c.Cy]
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
		IndentationLevel: 1,
		Chars:            []byte{},
		Length:           0,
		Highlighting:     []byte{},
		HlOpenComment:    false,
	}
}

func (b *Buffer) NewEditorAction(row Row, rowIndex int, actionType int, prevRowLength int, cx int, prevRow interface{}) *EditorAction {
	return &EditorAction{
		Row:        row,
		Index:      rowIndex,
		ActionType: actionType,
		PrevRow:    prevRow,
		Cx:         cx,
	}
}

func NewBuffer() *Buffer {
	return &Buffer{
		Rows:         []Row{},
		NumRows:      0,
		SearchState:  NewSearchState(),
		BufferSyntax: NewBufferSyntax(),
		UndoStack:    []EditorAction{},
		RedoStack:    []EditorAction{},
	}
}

func NewEditorConfig() *EditorConfig {
	return &EditorConfig{
		Cx:              0,
		Cy:              0,
		SliceIndex:      0,
		LineNumberWidth: 5,
		ScreenRows:      0,
		ScreenCols:      0,
		TerminalState:   nil,
		CurrentBuffer:   NewBuffer(),
		RowOff:          0,
		ColOff:          0,
		FileName:        "[Not Selected]",
		StatusMsg:       "",
		StatusMsgTime:   time.Time{},
		Dirty:           0,
		QuitTimes:       constants.QUIT_TIMES,
		Reader:          bufio.NewReader(os.Stdin),
		FirstRead:       true,
		UndoHistory:     30,
	}
}

func GetWindowSize(cfg *EditorConfig) error {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	cfg.ScreenCols = width
	cfg.ScreenRows = height
	return nil
}
