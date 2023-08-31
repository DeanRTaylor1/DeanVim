package constants

var TypeScriptSyntaxHighlighting map[string]byte = map[string]byte{
	// Control Flow
	"if":     CategoryConstants["controlFlow"],
	"else":   CategoryConstants["controlFlow"],
	"for":    CategoryConstants["controlFlow"],
	"switch": CategoryConstants["controlFlow"],
	"case":   CategoryConstants["controlFlow"],
	"break":  CategoryConstants["controlFlow"],
	"while":  CategoryConstants["controlFlow"],
	"do":     CategoryConstants["controlFlow"],
	"try":    CategoryConstants["controlFlow"],
	"catch":  CategoryConstants["controlFlow"],
	"throw":  CategoryConstants["exception"],

	// Variables and Constants
	"var": CategoryConstants["variable"],
	"let": CategoryConstants["variable"],

	"const": CategoryConstants["constant"],
	"env":   CategoryConstants["constant"],

	// Types
	"number":    CategoryConstants["type"],
	"string":    CategoryConstants["type"],
	"boolean":   CategoryConstants["type"],
	"enum":      CategoryConstants["type"],
	"any":       CategoryConstants["type"],
	"void":      CategoryConstants["type"],
	"null":      CategoryConstants["type"],
	"undefined": CategoryConstants["type"],

	// Functions
	"function":    CategoryConstants["function"],
	"return":      CategoryConstants["function"],
	"constructor": CategoryConstants["function"],

	// Operators
	"+": CategoryConstants["operator"],
	"-": CategoryConstants["operator"],
	"*": CategoryConstants["operator"],
	"/": CategoryConstants["operator"],

	// Keywords
	"this":       CategoryConstants["keyword"],
	"super":      CategoryConstants["keyword"],
	"interface":  CategoryConstants["keyword"],
	"class":      CategoryConstants["keyword"],
	"public":     CategoryConstants["keyword"],
	"private":    CategoryConstants["keyword"],
	"protected":  CategoryConstants["keyword"],
	"export":     CategoryConstants["keyword"],
	"import":     CategoryConstants["keyword"],
	"extends":    CategoryConstants["keyword"],
	"implements": CategoryConstants["keyword"],
	"instanceof": CategoryConstants["keyword"],
	"typeof":     CategoryConstants["keyword"],
	"as":         CategoryConstants["keyword"],
	"async":      CategoryConstants["keyword"],

	"process": CategoryConstants["builtin"],
	"console": CategoryConstants["builtin"],
	"from":    CategoryConstants["builtin"],

	// Boolean Literals
	"true":  CategoryConstants["boolean"],
	"false": CategoryConstants["boolean"],
}
