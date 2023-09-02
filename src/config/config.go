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

type EditorConfig struct {
	EditorMode             int
	Cx                     int
	Cy                     int
	SliceIndex             int
	LineNumberWidth        int
	ScreenRows             int
	ScreenCols             int
	TerminalState          *term.State
	CurrentBuffer          *Buffer
	Buffers                []Buffer
	RowOff                 int
	ColOff                 int
	FileName               string
	StatusMsg              string
	StatusMsgTime          time.Time
	QuitTimes              int
	Reader                 *bufio.Reader
	FirstRead              bool
	UndoHistory            int
	CurrentDirectory       string
	RootDirectory          string
	FileBrowserItems       []FileBrowserItem
	FileBrowserActionState FileBrowserActionState
	FileBrowserIntroLength int
	MotionBuffer           []rune
	MotionMap              map[string]func()
}

func (e *EditorConfig) ClearMotionBuffer() {
	e.MotionBuffer = []rune{}
}

func (e *EditorConfig) ExecuteMotion(motion string) bool {
	if action, exists := e.MotionMap[motion]; exists {
		action()
		return true
	}
	return false
}

func (e *EditorConfig) InstructionsLines() []string {
	return []string{
		"==========================================================",
		fmt.Sprintf("\x1b[36mGVim\x1b[39m Version: %s", constants.VERSION),
		"Gvim File Directory Preview",
		fmt.Sprintf("RootDir: %s", e.RootDirectory),
		"sorted by: name",
		"File Management: \x1b[38;5;1m%: Create R: Rename, D: Delete\x1b[39m",
		"Navigation:\x1b[38;5;1m<Up>/<k> to go up, <Down>/<j> to go down, <Enter> to select\x1b[39m",
		"==========================================================",
	}
}

func (e *EditorConfig) IsDir() bool {
	return e.CurrentSelectedFile().Name == ".." || e.CurrentSelectedFile().Type == "directory"
}

func (e *EditorConfig) CurrentSelectedFile() *FileBrowserItem {
	return &e.FileBrowserItems[e.Cy-len(e.InstructionsLines())]
}

// ReplaceBuffer replaces the buffer that matches the name of the current buffer with the current buffer's state.
func (e *EditorConfig) ReplaceBuffer() {
	for i, bufferItem := range e.Buffers {
		if bufferItem.Name == e.CurrentBuffer.Name {
			// Replace the item with the current state
			e.Buffers[i] = *e.CurrentBuffer
			break
		}
	}
}

// ReloadBuffer reloads the buffer that matches the name of the current buffer.
// It returns true if a matching buffer is found, false otherwise.
func (e *EditorConfig) ReloadBuffer(path string) bool {
	for _, bufferItem := range e.Buffers {
		if e.RootDirectory+bufferItem.Name == path {
			// Load the old buffer into CurrentBuffer
			e.CurrentBuffer = &bufferItem
			e.Cx = bufferItem.StoredCx
			e.Cy = bufferItem.StoredCy
			return true
		}
	}
	return false
}

func (e *EditorConfig) LoadNewBuffer() {
	e.Buffers = append(e.Buffers, *e.CurrentBuffer)
	e.CurrentBuffer.Idx = len(e.Buffers)
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
		QuitTimes:        constants.QUIT_TIMES,
		Reader:           bufio.NewReader(os.Stdin),
		FirstRead:        true,
		UndoHistory:      30,
		FileBrowserItems: []FileBrowserItem{},
		CurrentDirectory: "",
		MotionBuffer:     []rune{},
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
		e.CurrentBuffer.SliceIndex--
	}
}

func (e *EditorConfig) MoveCursorRight() {
	if e.EditorMode != constants.EDITOR_MODE_FILE_BROWSER {
		e.CurrentBuffer.SliceIndex++
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
