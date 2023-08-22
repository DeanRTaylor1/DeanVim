package highlighting

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/deanrtaylor1/go-editor/config"
	"github.com/deanrtaylor1/go-editor/constants"
)

func EditorSelectSyntaxHighlight(cfg *config.EditorConfig) {
	cfg.CurrentBuffer.BufferSyntax.FileType = "" // Reset to no filetype
	if cfg.FileName == "" {
		return
	}

	ext := filepath.Ext(cfg.FileName)
	for _, syntax := range cfg.CurrentBuffer.BufferSyntax.Syntaxes {
		for _, filematch := range syntax.FileMatch {
			isExt := strings.HasPrefix(filematch, ".")
			if (isExt && ext != "" && ext == filematch) || (!isExt && strings.Contains(cfg.FileName, filematch)) {
				cfg.CurrentBuffer.BufferSyntax.FileType = syntax.FileType
				cfg.CurrentBuffer.BufferSyntax.Flags = syntax.Flags
				cfg.CurrentBuffer.BufferSyntax.SingleLineCommentStart = syntax.SingleLineCommentStart
				cfg.CurrentBuffer.BufferSyntax.Keywords = syntax.Keywords
				for _, r := range cfg.CurrentBuffer.Rows {
					EditorUpdateSyntax(&r, cfg)
				}
				return
			}
		}
	}
}

func EditorUpdateSyntax(row *config.Row, cfg *config.EditorConfig) {
	row.Highlighting = make([]byte, row.Length)
	fill(row.Highlighting, constants.HL_NORMAL)

	if cfg.CurrentBuffer.BufferSyntax == nil {
		return
	}

	keywords := cfg.CurrentBuffer.BufferSyntax.Keywords

	scs := cfg.CurrentBuffer.BufferSyntax.SingleLineCommentStart
	scsLen := len(scs)

	prevSep := true
	inString := byte(0)
	i := 0
	for i < row.Length {
		c := row.Chars[i]
		prevHl := byte(constants.HL_NORMAL)
		if i > 0 {
			prevHl = row.Highlighting[i-1]
		}

		if scsLen > 0 && inString == 0 {
			if string(row.Chars[i:i+scsLen]) == scs {
				for ; i < row.Length; i++ {
					row.Highlighting[i] = constants.HL_COMMENT
				}
				break
			}
		}

		if cfg.CurrentBuffer.BufferSyntax.Flags&constants.HL_HIGHLIGHT_STRINGS != 0 {
			if inString != 0 {
				row.Highlighting[i] = constants.HL_STRING
				if c == '\\' && i+1 < row.Length {
					row.Highlighting[i+1] = constants.HL_STRING
					i += 2
					continue
				}
				if c == inString {
					inString = 0
				}
				i++
				prevSep = true
				continue
			} else if c == '"' || c == '\'' {
				inString = c
				row.Highlighting[i] = constants.HL_STRING
				i++
				continue
			}
		}

		if cfg.CurrentBuffer.BufferSyntax.Flags&constants.HL_HIGHLIGHT_NUMBERS != 0 {
			if unicode.IsDigit(rune(c)) && (prevSep || prevHl == constants.HL_NUMBER) || (c == '.' && prevHl == constants.HL_NUMBER) {
				row.Highlighting[i] = constants.HL_NUMBER
				i++
				prevSep = false
				continue
			}
		}

		if prevSep {
			foundKeyword := false
			for _, keyword := range keywords {
				klen := len(keyword)
				kw2 := keyword[klen-1] == '|'
				if kw2 {
					klen--
				}
				if i+klen <= len(row.Chars) && string(row.Chars[i:i+klen]) == keyword[:klen] && isSeparator(rune(row.Chars[i+klen])) {
					for k := 0; k < klen; k++ {
						row.Highlighting[i+k] = constants.HL_KEYWORD1
						if kw2 {
							row.Highlighting[i+k] = constants.HL_KEYWORD2
						}
					}
					i += klen
					foundKeyword = true
					break
				}
			}
			if foundKeyword {
				prevSep = false
				continue
			}
		}

		prevSep = isSeparator(rune(c))
		i++
	}
}

func isSeparator(c rune) bool {
	return unicode.IsSpace(c) || c == '\x00' || strings.ContainsRune(",.()+-/*=~%<>[];", c)
}

func fill(slice []byte, value byte) {
	for i := range slice {
		slice[i] = value
	}
}

func EditorSyntaxToColor(highlight byte) byte {
	switch highlight {
	case constants.HL_KEYWORD1:
		return 35
	case constants.HL_KEYWORD2:
		return 32
	case constants.HL_COMMENT:
		return 36
	case constants.HL_STRING:
		return 92
	case constants.HL_NUMBER:
		return 31
	case constants.HL_MATCH:
		return 34
	default:
		return 37
	}
}
