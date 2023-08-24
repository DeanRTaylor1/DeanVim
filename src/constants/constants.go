package constants

/** CONSTS **/
const VERSION = "0.0.1"

type SyntaxState int

const (
	STATE_NORMAL SyntaxState = iota
	STATE_SLCOMMENT
	STATE_MLCOMMENT
	STATE_STRING
	STATE_KEYWORD
	STATE_NUMBER
)

const (
	ARROW_LEFT  rune = 1000
	ARROW_RIGHT rune = 1001
	ARROW_UP    rune = 1002
	ARROW_DOWN  rune = 1003
	PAGE_UP     rune = 1004
	PAGE_DOWN   rune = 1005
	HOME_KEY    rune = 1006
	END_KEY     rune = 1007
	DEL_KEY     rune = 1008
	TAB_KEY     rune = 1009
	BACKSPACE   rune = 127
	QUIT_TIMES  int  = 3
	QUIT_KEY    rune = 'q'
	SAVE_KEY    rune = 's'
	ENTER_KEY   rune = '\r'
	TILDE       rune = '~'
	SPACE_RUNE  rune = ' '
)

const (
	ESCAPE_RESET_ATTRIBUTES  = "\x1b[m"
	ESCAPE_NEW_LINE          = "\r\n"
	ESCAPE_CLEAR_TO_LINE_END = "\x1b[K"
	ESCAPE_HIDE_CURSOR       = "\x1b[?25l"
	// MOVE CURSOR TO TOP LEFT
	ESCAPE_MOVE_TO_HOME_POS = "\x1b[H"

	ESCAPE_MOVE_TO_COORDS = "\x1b[%d;%dH"
	ESCAPE_SHOW_CURSOR    = "\x1b[?25h"
	ESCAPE_CLEAR_SCREEN   = "\033[2J"
)

const (
	TEXT_BLACK          = "\x1b[30m"
	TEXT_RED            = "\x1b[31m"
	TEXT_GREEN          = "\x1b[32m"
	TEXT_YELLOW         = "\x1b[33m"
	TEXT_BLUE           = "\x1b[34m"
	TEXT_MAGENTA        = "\x1b[35m"
	TEXT_CYAN           = "\x1b[36m"
	TEXT_WHITE          = "\x1b[37m"
	TEXT_BRIGHT_BLACK   = "\x1b[90m"
	TEXT_BRIGHT_RED     = "\x1b[91m"
	TEXT_BRIGHT_GREEN   = "\x1b[92m"
	TEXT_BRIGHT_YELLOW  = "\x1b[93m"
	TEXT_BRIGHT_BLUE    = "\x1b[94m"
	TEXT_BRIGHT_MAGENTA = "\x1b[95m"
	TEXT_BRIGHT_CYAN    = "\x1b[96m"
	TEXT_BRIGHT_WHITE   = "\x1b[97m"
	BACKGROUND_BLACK    = "\x1b[40m"
	BACKGROUND_RED      = "\x1b[41m"
	BACKGROUND_GREEN    = "\x1b[42m"
	BACKGROUND_YELLOW   = "\x1b[43m"
	BACKGROUND_BLUE     = "\x1b[44m"
	BACKGROUND_MAGENTA  = "\x1b[45m"
	BACKGROUND_CYAN     = "\x1b[46m"
	BACKGROUND_WHITE    = "\x1b[47m"
	RESET               = "\x1b[0m"
	BOLD                = "\x1b[1m"
	UNDERLINE           = "\x1b[4m"
	FOREGROUND_RESET    = "\x1b[39m"
	BACKGROUND_RESET    = "\x1b[49m"
)

const TAB_STOP = 4

const (
	HL_NORMAL = iota
	HL_NUMBER
	HL_MATCH
	HL_STRING
	HL_COMMENT
	HL_MLCOMMENT
	HL_KEYWORD1
	HL_KEYWORD2
	HL_CONTROL_FLOW  // for "if", "else", "for", "while", "switch"
	HL_VARIABLE      // for "var", "let"
	HL_CONSTANT      // for "const", "enum"
	HL_TYPE          // for "int", "float", "string", "bool"
	HL_FUNCTION      // for "func", "return"
	HL_PREPROCESSOR  // for "#include", "#define"
	HL_STORAGE_CLASS // for "static", "extern", "public", "private"
	HL_OPERATOR      // for "+", "-", "*", "/"
	HL_BOOLEAN       // for "true", "false"
	HL_KEYWORD       // for general language keywords
	HL_BUILTIN       // for built-in functions or types
	HL_ANNOTATION    // for annotations or attributes
	HL_EXCEPTION     // for "throw", "try", "catch"
	HL_MODULE        // for "import", "package", "module"
	HL_DEBUG         // for debug-related keywords
	HL_TEST          // for test-related keywords
	HL_DOCUMENTATION // for documentation comments
)

const (
	HL_HIGHLIGHT_NUMBERS = 1 << iota
	HL_HIGHLIGHT_STRINGS
)
