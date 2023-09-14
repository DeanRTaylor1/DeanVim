package grep

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunGrep(searchTerm string, directory string) ([]string, error) {
	var results []string

	// Run the grep command
	cmd := exec.Command("grep", "-r", "-n", searchTerm, directory)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Parse the output
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Split the line into filename, line number, and content
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}

		filename := parts[0]
		lineNumber := parts[1]
		content := parts[2]
		content = strings.ReplaceAll(content, "\t", "    ")

		// Remove the directory prefix from the filename
		filename = strings.TrimPrefix(filename, directory)
		if strings.HasPrefix(filename, "/") {
			filename = strings.TrimPrefix(filename, "/")
		}

		// Find the character position
		charPosition := strings.Index(content, searchTerm)
		if charPosition == -1 {
			continue
		}

		// Format the result to include the full line content
		result := fmt.Sprintf("%s:%s:%d:%s", filename, lineNumber, charPosition, content)
		results = append(results, result)
	}

	return results, nil
}
