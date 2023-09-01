package actions

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
	"github.com/deanrtaylor1/go-editor/highlighting"
)

func ReadHandler(cfg *config.EditorConfig, arg string) {
	fileInfo, err := os.Stat(arg)
	if err != nil {
		log.Fatal(err)
	}

	if arg == "." {
		cfg.SetMode(constants.EDITOR_MODE_FILE_BROWSER)
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not get current directory")
		}
		if cfg.RootDirectory == "" {
			cfg.RootDirectory = currentDir
		}
		DirectoryOpen(cfg, currentDir)
	} else if fileInfo.IsDir() {
		cfg.SetMode(constants.EDITOR_MODE_FILE_BROWSER)
		// Set the current directory path in the config
		if cfg.RootDirectory == "" {
			cfg.RootDirectory = arg
		}
		DirectoryOpen(cfg, arg)
	} else {
		cfg.EditorMode = constants.EDITOR_MODE_NORMAL
		if cfg.CurrentBuffer.Name != "" {
			cfg.ReplaceBuffer()
		}
		foundBuffer := cfg.ReloadBuffer(arg)
		if !foundBuffer {
			err := EditorOpen(cfg, arg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func DirectoryOpen(cfg *config.EditorConfig, path string) error {
	// Read the directory
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Clear the existing FileBrowserItems
	cfg.FileBrowserItems = []config.FileBrowserItem{}

	// Populate the FileBrowserItems slice
	for _, entry := range dirEntries {

		fullPath := filepath.Join(path, entry.Name())
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			return err
		}

		item := config.FileBrowserItem{
			Name:       entry.Name(),
			Path:       fullPath,
			CreatedAt:  fileInfo.ModTime(),
			ModifiedAt: fileInfo.ModTime(),
		}

		if entry.Type().IsDir() {
			item.Type = "directory"
			item.Extension = "directory" // or leave it empty
		} else {
			item.Type = "file"
			ext := filepath.Ext(entry.Name()) // Remove the leading dot
			if len(ext) > 1 {
				item.Extension = ext[1:]
			}
		}

		cfg.FileBrowserItems = append(cfg.FileBrowserItems, item)
	}

	// Sort the FileBrowserItems so that directories appear first
	sort.Slice(cfg.FileBrowserItems, func(i, j int) bool {
		return cfg.FileBrowserItems[i].Type == "directory" && cfg.FileBrowserItems[j].Type != "directory"
	})

	if path != "/" {
		parentDir := filepath.Dir(path)
		parentItem := config.FileBrowserItem{
			Name:       "..",
			Path:       parentDir,
			Type:       "directory",
			Extension:  "directory",
			CreatedAt:  time.Time{},
			ModifiedAt: time.Time{},
		}
		cfg.FileBrowserItems = append([]config.FileBrowserItem{parentItem}, cfg.FileBrowserItems...)
	}

	cfg.CurrentDirectory = path

	EditorSetStatusMessage(cfg, fmt.Sprintf("%s", cfg.RootDirectory))

	return nil
}

func EditorOpen(cfg *config.EditorConfig, fileName string) error {
	cfg.EditorMode = constants.EDITOR_MODE_NORMAL
	if !cfg.FirstRead {
		cfg.CurrentBuffer = config.NewBuffer()
	}
	cfg.Cx = 0
	cfg.Cy = 0
	cfg.CurrentBuffer.SliceIndex = 0

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening file")
	}
	defer file.Close()
	relativeFileName := strings.TrimPrefix(fileName, cfg.RootDirectory)

	// If the fileName didn't start with RootDirectory, just use the base name
	if relativeFileName == fileName {
		relativeFileName = filepath.Base(fileName)
	}

	cfg.FileName = relativeFileName

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
		EditorInsertRow(row, row.Idx, cfg)
		cfg.CurrentBuffer.NumRows++ // Update NumRows within CurrentBuffer
	}
	highlighting.HighlightFileFromRow(0, cfg)

	if err := scanner.Err(); err != nil {
		return err
	}
	cfg.CurrentBuffer.Dirty = 0
	cfg.FirstRead = false
	cfg.CurrentBuffer.Name = relativeFileName
	if len(cfg.Buffers) < 1 {
		cfg.Buffers = make([]config.Buffer, 10)
	}

	cfg.LoadNewBuffer()

	EditorSetStatusMessage(cfg, "HELP: CTRL-S = Save | Ctrl-Q = quit | Ctr-f = find")

	return nil
}

func EditorSave(cfg *config.EditorConfig) (string, error) {
	if cfg.CurrentBuffer.Name == "" {
		return "", errors.New("no filename provided")
	}

	startTime := time.Now()
	content := EditorRowsToString(cfg)

	file, err := os.OpenFile(fmt.Sprintf("%s%s", cfg.RootDirectory, cfg.CurrentBuffer.Name), os.O_RDWR|os.O_CREATE, 0644)
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
	message := fmt.Sprintf("\"%s\", %dL, %dB, %.3fms: written", cfg.CurrentBuffer.Name, numLines, numBytes, float64(elapsedTime.Nanoseconds())/1e6)

	cfg.CurrentBuffer.Dirty = 0

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
