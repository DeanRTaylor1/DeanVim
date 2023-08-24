package constants

var JavaScriptSyntaxHighlighting map[string]byte = map[string]byte{
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
	"var":   CategoryConstants["variable"],
	"let":   CategoryConstants["variable"],
	"const": CategoryConstants["constant"],

	// Functions
	"function": CategoryConstants["function"],
	"return":   CategoryConstants["function"],

	// Operators
	"+": CategoryConstants["operator"],
	"-": CategoryConstants["operator"],
	"*": CategoryConstants["operator"],
	"/": CategoryConstants["operator"],

	// Keywords
	"this":       CategoryConstants["keyword"],
	"super":      CategoryConstants["keyword"],
	"class":      CategoryConstants["keyword"],
	"export":     CategoryConstants["keyword"],
	"import":     CategoryConstants["keyword"],
	"extends":    CategoryConstants["keyword"],
	"instanceof": CategoryConstants["keyword"],
	"typeof":     CategoryConstants["keyword"],
	"new":        CategoryConstants["keyword"],
	"delete":     CategoryConstants["keyword"],
	"in":         CategoryConstants["keyword"],
	"of":         CategoryConstants["keyword"],

	// Boolean Literals
	"true":  CategoryConstants["boolean"],
	"false": CategoryConstants["boolean"],
}
