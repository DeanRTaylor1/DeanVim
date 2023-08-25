package actions

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func EditorOpen(cfg *config.EditorConfig, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening file")
	}
	defer file.Close()
	cfg.FileName = file.Name()

	highlighting.EditorSelectSyntaxHighlight(cfg)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		linelen := len(line)
		for linelen > 0 && (line[linelen-1] == '\n' || line[linelen-1] == '\r') {
			linelen--
		}
		row := config.NewRow() // Create a new Row using the NewRow function
		row.Chars = []byte(line[:linelen])
		row.Length = linelen
		row.Idx = len(cfg.CurrentBuffer.Rows)
		row.Highlighting = make([]byte, row.Length)
		highlighting.Fill(row.Highlighting, constants.HL_NORMAL)
		row.Tabs = make([]byte, row.Length)
		highlighting.SyntaxHighlightStateMachine(row, cfg)
		EditorInsertRow(row, -1, cfg)
		cfg.CurrentBuffer.NumRows++ // Update NumRows within CurrentBuffer
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	cfg.Dirty = 0
	cfg.FirstRead = false

	return nil
}

func EditorSave(cfg *config.EditorConfig) (string, error) {
	if cfg.FileName == "[Not Selected]" {
		return "", errors.New("no filename provided")
	}

	startTime := time.Now()
	content := EditorRowsToString(cfg)

	file, err := os.OpenFile(cfg.FileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := file.Truncate(int64(len(content))); err != nil {
		return "", fmt.Errorf("failed to truncate file: %w", err)
	}

	n, err := file.WriteString(content)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}
	if n != len(content) {
		return "", errors.New("unexpected number of bytes written to file")
	}

	elapsedTime := time.Since(startTime) // End timing
	numLines := len(cfg.CurrentBuffer.Rows)
	numBytes := len(content)
	message := fmt.Sprintf("\"%s\", %dL, %dB, %.3fms: written", cfg.FileName, numLines, numBytes, float64(elapsedTime.Nanoseconds())/1e6)

	cfg.Dirty = 0

	return message, nil
}

func EditorRowsToString(cfg *config.EditorConfig) string {
	var buffer strings.Builder
	for _, row := range cfg.CurrentBuffer.Rows {
		buffer.Write(row.Chars)
		buffer.WriteByte('\n')
	}
	return buffer.String()
}
