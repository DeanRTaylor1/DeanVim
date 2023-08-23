package treesitter

import (
	"errors"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
	"github.com/smacker/go-tree-sitter/yaml"
)

func GetParserForFileType(fileType string) (*sitter.Parser, error) {
	var lang *sitter.Language
	switch fileType {
	case "go":
		lang = javascript.GetLanguage()
		break
	case "typescript":
		lang = typescript.GetLanguage()
		break
	case "javascript":
		lang = javascript.GetLanguage()
		break
	case "yaml":
		lang = yaml.GetLanguage()
		break
	case "json":
		lang = javascript.GetLanguage()
	default:
		return nil, errors.New("Invalid file type")
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)
	return parser, nil
}
