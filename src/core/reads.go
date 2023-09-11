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
func ReadHandler(e *config.Editor, arg string) {
	fileInfo, err := os.Stat(arg)
	if err != nil {
		log.Fatal(err)
	}

	if arg == "." {
		e.SetMode(constants.EDITOR_MODE_FILE_BROWSER)
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not get current directory")
		}
		if e.RootDirectory == "" {
			e.RootDirectory = currentDir
		}
		DirectoryOpen(e, currentDir)
	} else if fileInfo.IsDir() {
		e.SetMode(constants.EDITOR_MODE_FILE_BROWSER)
		// Set the current directory path in the config
		if e.RootDirectory == "" {
			e.RootDirectory = arg
		}
		DirectoryOpen(e, arg)
	} else {
		if e.RootDirectory == "" {
			currentDir, err := os.Getwd()
			if err != nil {
				log.Fatal("Could Not Get the current directory")
			}
			e.RootDirectory = currentDir

		}
		e.EditorMode = constants.EDITOR_MODE_NORMAL
		if e.CurrentBuffer.Name != "" {
			e.ReplaceBuffer()
		}
		foundBuffer := e.ReloadBuffer(arg)
		if !foundBuffer {
			err := FileOpen(e, arg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func DirectoryOpen(e *config.Editor, path string) error {
	// Read the directory
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Clear the existing FileBrowserItems
	e.FileBrowserItems = []config.FileBrowserItem{}

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

		e.FileBrowserItems = append(e.FileBrowserItems, item)
	}

	// Sort the FileBrowserItems so that directories appear first
	sort.Slice(e.FileBrowserItems, func(i, j int) bool {
		return e.FileBrowserItems[i].Type == "directory" && e.FileBrowserItems[j].Type != "directory"
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
		e.FileBrowserItems = append([]config.FileBrowserItem{parentItem}, e.FileBrowserItems...)
	}

	e.CurrentDirectory = path

	// EditorSetStatusMessage(e, fmt.Sprintf("%s", e.RootDirectory))

	return nil
}
