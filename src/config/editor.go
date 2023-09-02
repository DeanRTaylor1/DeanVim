package config

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/utils"
	"golang.org/x/term"
)

type Editor struct {
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

func NewEditor() *Editor {
	return &Editor{
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

func (e *Editor) ResetCursorCoords() {
	e.Cx = 0
	e.Cy = 0
}

func (e *Editor) CacheCursorCoords() {

	e.CurrentBuffer.StoredCx = e.Cx
	e.CurrentBuffer.StoredCy = e.Cy
}

func (e *Editor) ClearMotionBuffer() {
	e.MotionBuffer = []rune{}
}

func (e *Editor) ExecuteMotion(motion string) bool {
	if action, exists := e.MotionMap[motion]; exists {
		action()
		return true
	}
	return false
}

func (e *Editor) InstructionsLines() []string {
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

func (e *Editor) IsDir() bool {
	return e.CurrentSelectedFile().Name == ".." || e.CurrentSelectedFile().Type == "directory"
}

func (e *Editor) CurrentSelectedFile() *FileBrowserItem {
	return &e.FileBrowserItems[e.Cy-len(e.InstructionsLines())]
}

// ReplaceBuffer replaces the buffer that matches the name of the current buffer with the current buffer's state.
func (e *Editor) ReplaceBuffer() {
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
func (e *Editor) ReloadBuffer(path string) bool {
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

func (e *Editor) LoadNewBuffer() {
	e.Buffers = append(e.Buffers, *e.CurrentBuffer)
	e.CurrentBuffer.Idx = len(e.Buffers)
}

func (c *Editor) ClearRedoStack() {
	c.CurrentBuffer.RedoStack = []EditorAction{}
}

func (c *Editor) GetAdjustedCx() int {
	adjustedCx := utils.Max(0, c.Cx-c.LineNumberWidth)
	adjustedCx = utils.Min(adjustedCx, len(c.CurrentBuffer.Rows[c.Cy].Chars))
	return adjustedCx
}

func (c *Editor) GetCurrentRow() *Row {
	return &c.CurrentBuffer.Rows[c.Cy]
}

func (e *Editor) SpecialRefreshCase() bool {
	return e.Cx >= e.ScreenCols-e.LineNumberWidth || e.Cy >= e.ScreenRows || e.Cx-e.LineNumberWidth < e.ColOff || e.Cy-e.RowOff < 0 || e.Cx-e.ColOff == 5
}

func (e *Editor) IsBrowsingFiles() bool {
	return e.EditorMode == constants.EDITOR_MODE_FILE_BROWSER
}

func (e *Editor) SetMode(mode int) {
	e.EditorMode = mode
}

func (e *Editor) MoveCursorLeft() {
	e.Cx--
	if e.EditorMode != constants.EDITOR_MODE_FILE_BROWSER {
		e.CurrentBuffer.SliceIndex--
	}
}

func (e *Editor) MoveCursorRight() {
	if e.EditorMode != constants.EDITOR_MODE_FILE_BROWSER {
		e.CurrentBuffer.SliceIndex++
	}
	e.Cx++
}

func (e *Editor) MoveCursorUp() {
	e.Cy--
}

func (e *Editor) MoveCursorDown() {
	e.Cy++
}
