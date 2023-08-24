package constants

var GoSyntaxHighlighting map[string]byte = map[string]byte{
	// Control Flow
	"if":     CategoryConstants["controlFlow"],
	"else":   CategoryConstants["controlFlow"],
	"for":    CategoryConstants["controlFlow"],
	"switch": CategoryConstants["controlFlow"],
	"case":   CategoryConstants["controlFlow"],

	// Variables and Constants
	"var":   CategoryConstants["variable"],
	"const": CategoryConstants["constant"],

	// Module
	"import": CategoryConstants["module"],

	// Types
	"int":        CategoryConstants["type"],
	"int8":       CategoryConstants["type"],
	"int16":      CategoryConstants["type"],
	"int32":      CategoryConstants["type"],
	"int64":      CategoryConstants["type"],
	"uint":       CategoryConstants["type"],
	"uint8":      CategoryConstants["type"],
	"uint16":     CategoryConstants["type"],
	"uint32":     CategoryConstants["type"],
	"uint64":     CategoryConstants["type"],
	"uintptr":    CategoryConstants["type"],
	"float32":    CategoryConstants["type"],
	"float64":    CategoryConstants["type"],
	"complex64":  CategoryConstants["type"],
	"complex128": CategoryConstants["type"],
	"string":     CategoryConstants["type"],
	"bool":       CategoryConstants["type"],
	"byte":       CategoryConstants["type"],
	"rune":       CategoryConstants["type"],
	"error":      CategoryConstants["type"],

	// Functions
	"func":   CategoryConstants["function"],
	"return": CategoryConstants["function"],

	// Operators
	"+": CategoryConstants["operator"],
	"-": CategoryConstants["operator"],
	"*": CategoryConstants["operator"],
	"/": CategoryConstants["operator"],

	// Builtins
	"append": CategoryConstants["builtin"],
	"len":    CategoryConstants["builtin"],
	"make":   CategoryConstants["builtin"],
	"new":    CategoryConstants["builtin"],
	"cap":    CategoryConstants["builtin"],
	"close":  CategoryConstants["builtin"],
	"copy":   CategoryConstants["builtin"],
	"delete": CategoryConstants["builtin"],

	"true":  CategoryConstants["boolean"],
	"false": CategoryConstants["boolean"],
	"goto":  CategoryConstants["controlFlow"],
	"break": CategoryConstants["controlFlow"],
}
