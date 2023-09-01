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
	ActionType   int
	Row          Row
	Index        int
	PrevRow      interface{}
	Cx           int
	RedoFunction func()
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

type EditorConfig struct {
	EditorMode       int
	Cx               int
	Cy               int
	SliceIndex       int
	LineNumberWidth  int
	ScreenRows       int
	ScreenCols       int
	TerminalState    *term.State
	CurrentBuffer    *Buffer
	RowOff           int
	ColOff           int
	FileName         string
	StatusMsg        string
	StatusMsgTime    time.Time
	Dirty            int
	QuitTimes        int
	Reader           *bufio.Reader
	FirstRead        bool
	UndoHistory      int
	CurrentDirectory string
	FileBrowserItems []FileBrowserItem
}

func (c *EditorConfig) ClearRedoStack() {
	c.CurrentBuffer.RedoStack = []EditorAction{}
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

func NewEditorConfig() *EditorConfig {
	return &EditorConfig{
		EditorMode:       constants.EDITOR_MODE_NORMAL,
		Cx:               0,
		Cy:               0,
		SliceIndex:       0,
		LineNumberWidth:  5,
		ScreenRows:       0,
		ScreenCols:       0,
		TerminalState:    nil,
		CurrentBuffer:    NewBuffer(),
		RowOff:           0,
		ColOff:           0,
		FileName:         "[Not Selected]",
		StatusMsg:        "",
		StatusMsgTime:    time.Time{},
		Dirty:            0,
		QuitTimes:        constants.QUIT_TIMES,
		Reader:           bufio.NewReader(os.Stdin),
		FirstRead:        true,
		UndoHistory:      30,
		FileBrowserItems: []FileBrowserItem{},
		CurrentDirectory: "",
	}
}

func (e *EditorConfig) SpecialRefreshCase() bool {
	return e.Cx >= e.ScreenCols-e.LineNumberWidth || e.Cy >= e.ScreenRows || e.Cx-e.LineNumberWidth < e.ColOff || e.Cy-e.RowOff < 0 || e.Cx-e.ColOff == 5
}

func (e *EditorConfig) IsBrowsingFiles() bool {
	return e.EditorMode == constants.EDITOR_MODE_FILE_BROWSER
}

func (e *EditorConfig) SetMode(mode int) {
	e.EditorMode = mode
}

func (e *EditorConfig) MoveCursorLeft() {
	e.Cx--
	if e.EditorMode != constants.EDITOR_MODE_FILE_BROWSER {
		e.SliceIndex--
	}
}

func (e *EditorConfig) MoveCursorRight() {
	if e.EditorMode != constants.EDITOR_MODE_FILE_BROWSER {
		e.SliceIndex++
	}
	e.Cx++
}

func (e *EditorConfig) MoveCursorUp() {
	e.Cy--
}

func (e *EditorConfig) MoveCursorDown() {
	e.Cy++
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
