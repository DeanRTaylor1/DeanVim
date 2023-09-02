package core

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

// this function checks the type of the item and directs to the relevant function
func ReadHandler(cfg *config.Editor, arg string) {
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
			err := FileOpen(cfg, arg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func DirectoryOpen(cfg *config.Editor, path string) error {
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

	// EditorSetStatusMessage(cfg, fmt.Sprintf("%s", cfg.RootDirectory))

	return nil
}
