package constants

/** CONSTS **/
const VERSION = "0.0.1"

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
	BACKSPACE   rune = 127
	QUIT_TIMES  int  = 3
	QUIT_KEY    rune = 'q'
	SAVE_KEY    rune = 's'
	ENTER_KEY   rune = '\r'
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
)

const TAB_STOP = 4

const (
	HL_NORMAL = iota
	HL_NUMBER
)
