package constants

var CategoryConstants map[string]byte = map[string]byte{
	"controlFlow":   HL_CONTROL_FLOW,  // for "if", "else", "for", "while", "switch"
	"variable":      HL_VARIABLE,      // for "var", "let"
	"constant":      HL_CONSTANT,      // for "const", "enum"
	"type":          HL_TYPE,          // for "int", "float", "string", "bool"
	"function":      HL_FUNCTION,      // for "func", "return"
	"preprocessor":  HL_PREPROCESSOR,  // for "#include", "#define"
	"storageClass":  HL_STORAGE_CLASS, // for "static", "extern", "public", "private"
	"operator":      HL_OPERATOR,      // for "+", "-", "*", "/"
	"comment":       HL_COMMENT,       // for "//", "/*", "*/"
	"string":        HL_STRING,        // for string literals
	"number":        HL_NUMBER,        // for numeric literals
	"boolean":       HL_BOOLEAN,       // for "true", "false"
	"keyword":       HL_KEYWORD,       // for general language keywords
	"builtin":       HL_BUILTIN,       // for built-in functions or types
	"annotation":    HL_ANNOTATION,    // for annotations or attributes
	"exception":     HL_EXCEPTION,     // for "throw", "try", "catch"
	"module":        HL_MODULE,        // for "import", "package", "module"
	"debug":         HL_DEBUG,         // for debug-related keywords
	"test":          HL_TEST,          // for test-related keywords
	"documentation": HL_DOCUMENTATION, // for documentation comments
}

type SyntaxHighlighting struct {
	FileType               string
	FileMatch              []string
	SingleLineCommentStart string
	MultiLineCommentStart  string
	MultiLineCommentEnd    string
	Flags                  int
	Syntax                 map[string]byte
}

var Syntaxes = []SyntaxHighlighting{
	{
		FileType:               "go",
		FileMatch:              []string{".go"},
		SingleLineCommentStart: "//",
		MultiLineCommentStart:  "/*",
		MultiLineCommentEnd:    "*/",
		Flags:                  HL_HIGHLIGHT_NUMBERS | HL_HIGHLIGHT_STRINGS,
		Syntax:                 GoSyntaxHighlighting,
	},
	{
		FileType:               "typescript",
		FileMatch:              []string{".ts", ".tsx"},
		SingleLineCommentStart: "//",
		MultiLineCommentStart:  "/*",
		MultiLineCommentEnd:    "*/",
		Flags:                  HL_HIGHLIGHT_NUMBERS | HL_HIGHLIGHT_STRINGS,
		Syntax:                 TypeScriptSyntaxHighlighting,
	},
	// {
	// 	FileType:               "rust",
	// 	FileMatch:              []string{".rs"},
	// 	SingleLineCommentStart: "//",
	// 	MultiLineCommentStart:  "/*",
	// 	MultiLineCommentEnd:    "*/",
	// 	Flags:                  HL_HIGHLIGHT_NUMBERS | HL_HIGHLIGHT_STRINGS,
	// 	Keywords:               []string{"fn", "let", "const", "trait", "struct", "enum", "return", "if", "else"},
	// },
	{
		FileType:               "javascript",
		FileMatch:              []string{".js", ".jsx"},
		SingleLineCommentStart: "//",
		MultiLineCommentStart:  "/*",
		MultiLineCommentEnd:    "*/",
		Flags:                  HL_HIGHLIGHT_NUMBERS | HL_HIGHLIGHT_STRINGS,
		Syntax:                 JavaScriptSyntaxHighlighting,
	},
}

var BracketPairs = map[rune]rune{
	'{': '}',
	'[': ']',
	'(': ')',
}
