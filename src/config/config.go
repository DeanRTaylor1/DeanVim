package config

import (
	"bufio"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/deanrtaylor1/go-editor/constants"
)

type buffer struct {
	Rows    []Row
	NumRows int
}

type Row struct {
	Row          []byte
	length       int
	Highlighting []byte
}

type EditorConfig struct {
	Cx            int
	Cy            int
	ScreenRows    int
	ScreenCols    int
	TerminalState *term.State
	NumRows       int
	Rows          [][]byte
	RowOff        int
	ColOff        int
	FileName      string
	StatusMsg     string
	StatusMsgTime time.Time
	Dirty         int
	QuitTimes     int
	Reader        *bufio.Reader
	Highlighting  [][]byte
}

func NewEditorConfig() *EditorConfig {
	return &EditorConfig{
		Cx:            0,
		Cy:            0,
		ScreenRows:    0,
		ScreenCols:    0,
		TerminalState: nil,
		Rows:          [][]byte{},
		NumRows:       0,
		RowOff:        0,
		ColOff:        0,
		FileName:      "[Not Selected]",
		StatusMsg:     "",
		StatusMsgTime: time.Time{},
		Dirty:         0,
		QuitTimes:     constants.QUIT_TIMES,
		Reader:        bufio.NewReader(os.Stdin),
		Highlighting:  [][]byte{},
	}
}
