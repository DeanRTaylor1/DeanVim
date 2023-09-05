package config

import (
	"github.com/deanrtaylor1/go-editor/constants"
)

type Buffer struct {
	Idx                int
	Name               string
	Rows               []Row
	NumRows            int
	SearchState        *SearchState
	BufferSyntax       *BufferSyntax
	UndoStack          []EditorAction
	RedoStack          []EditorAction
	NeedsFullHighlight bool
	SliceIndex         int
	Dirty              int
	StoredCx           int
	StoredCy           int
	SelectedCxStart    int
	SelectedCyStart    int
	SelectedCyEnd      int
	SelectedCxEnd      int
}

type BufferSyntax struct {
	FileType               string
	Flags                  int
	SingleLineCommentStart string
	MultiLineCommentStart  string
	MultiLineCommentEnd    string
	Keywords               map[string]byte
	Syntaxes               []constants.SyntaxHighlighting
}

func (b *Buffer) MultiLineHighlight() bool {
	return b.SelectedCyStart != b.SelectedCyEnd
}

func (b *Buffer) ClearSelection() {
	b.SelectedCxEnd = 0
	b.SelectedCyEnd = 0
	b.SelectedCxStart = 0
	b.SelectedCyStart = 0
}

func (b *Buffer) ReplaceRowAtIndex(index int, newRow Row) {
	if index < 0 || index >= len(b.Rows) {
		return
	}

	b.Rows[index] = newRow
}

func (b *Buffer) RemoveRowAtIndex(index int) {
	if index < 0 || index >= len(b.Rows) {
		return
	}

	beforeRows := b.Rows[:index]

	afterRows := b.Rows[index+1:]

	b.Rows = append(beforeRows, afterRows...)
	b.NumRows--
}

func (b *Buffer) InsertRowAtIndex(index int, newRow Row) {
	beforeRows := b.Rows[:index]

	newRows := []Row{newRow}

	afterRows := b.Rows[index:]

	b.Rows = append(beforeRows, append(newRows, afterRows...)...)
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

func (b *Buffer) NewEditorAction(row Row, rowIndex int, actionType int, prevRowLength int, cx int, prevRow interface{}, redoFunction func()) *EditorAction {
	return &EditorAction{
		Row:          row,
		Index:        rowIndex,
		ActionType:   actionType,
		PrevRow:      prevRow,
		Cx:           cx,
		RedoFunction: redoFunction,
	}
}

func NewBuffer() *Buffer {
	return &Buffer{
		Rows:               []Row{},
		NumRows:            0,
		SearchState:        NewSearchState(),
		BufferSyntax:       NewBufferSyntax(),
		UndoStack:          []EditorAction{},
		RedoStack:          []EditorAction{},
		NeedsFullHighlight: false,
		SliceIndex:         0,
		Dirty:              0,
	}
}
