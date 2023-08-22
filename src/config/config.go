package config

import (
	"bufio"
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

type SearchState struct {
	LastMatch   int
	Direction   int
	SavedHlLine int
	SavedHl     []byte
}

type Buffer struct {
	Rows        []Row
	NumRows     int
	SearchState *SearchState
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

func NewSearchState() *SearchState {
	return &SearchState{
		LastMatch:   -1,
		Direction:   1,
		SavedHlLine: 0,
		SavedHl:     []byte{},
	}
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
		Rows:        []Row{},
		NumRows:     0,
		SearchState: NewSearchState(),
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
