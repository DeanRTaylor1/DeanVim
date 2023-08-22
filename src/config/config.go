package config

import (
	"bufio"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/deanrtaylor1/go-editor/constants"
)

type Buffer struct {
	Rows    []Row
	NumRows int
}

type Row struct {
	Chars        []byte
	Length       int
	Highlighting []byte
}

type EditorConfig struct {
	Cx            int
	Cy            int
	ScreenRows    int
	ScreenCols    int
	TerminalState *term.State
	// NumRows       int
	CurrentBuffer *Buffer
	RowOff        int
	ColOff        int
	FileName      string
	StatusMsg     string
	StatusMsgTime time.Time
	Dirty         int
	QuitTimes     int
	Reader        *bufio.Reader
	// Highlighting  [][]byte
}

func NewRow() *Row {
	return &Row{
		Chars:        []byte{},
		Length:       0,
		Highlighting: []byte{},
	}
}

func NewBuffer() *Buffer {
	return &Buffer{
		Rows:    []Row{},
		NumRows: 0,
	}
}

func NewEditorConfig() *EditorConfig {
	return &EditorConfig{
		Cx:            0,
		Cy:            0,
		ScreenRows:    0,
		ScreenCols:    0,
		TerminalState: nil,
		CurrentBuffer: NewBuffer(),
		// NumRows:       0,
		RowOff:        0,
		ColOff:        0,
		FileName:      "[Not Selected]",
		StatusMsg:     "",
		StatusMsgTime: time.Time{},
		Dirty:         0,
		QuitTimes:     constants.QUIT_TIMES,
		Reader:        bufio.NewReader(os.Stdin),
		// Highlighting:  [][]byte{},
	}
}
