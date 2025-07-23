package main

import (
	"strings"
	"unicode"
)

// highlight types
const (
	HL_NORMAL = iota + 0
	HL_COMMENT
	HL_MULTILINE_COMMENT
	HL_KEYWORD1
	HL_KEYWORD2
	HL_STRING
	HL_NUMBER
	HL_SEARCH_MATCH
)

// highlight flags
const (
	HL_HIGHLIGHT_NUMBERS = 1 << iota
	HL_HIGHLIGHT_STRINGS
)

const (
	BLACK = 30 + iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

type Syntax struct {
	fileType               []byte
	fileMatch              [][]byte
	keywords               [][]byte
	singleLineCommentStart []byte
	multiLineCommentStart  []byte
	multiLineCommentEnd    []byte
	flags                  int
}

type SyntaxDB struct {
	syntax []Syntax
}

var goSyntax = Syntax{
	fileType: []byte("Go"),
	fileMatch: [][]byte{
		[]byte(".go"),
	},
	keywords: [][]byte{
		[]byte("break"),
		[]byte("case"),
		[]byte("chan"),
		[]byte("const"),
		[]byte("continue"),
		[]byte("default"),
		[]byte("defer"),
		[]byte("do"),
		[]byte("else"),
		[]byte("fallthrough"),
		[]byte("for"),
		[]byte("func"),
		[]byte("go"),
		[]byte("if"),
		[]byte("import"),
		[]byte("interface"),
		[]byte("map"),
		[]byte("package"),
		[]byte("range"),
		[]byte("return"),
		[]byte("select"),
		[]byte("struct"),
		[]byte("switch"),
		[]byte("type"),
		[]byte("var"),
		[]byte("string|"),
		[]byte("bool|"),
		[]byte("int|"),
		[]byte("int8|"),
		[]byte("int16|"),
		[]byte("int32|"),
		[]byte("int32|"),
		[]byte("int64|"),
		[]byte("uint|"),
		[]byte("uint8|"),
		[]byte("uint16|"),
		[]byte("uint32|"),
		[]byte("uint64|"),
		[]byte("float32|"),
		[]byte("float64|"),
		[]byte("complex64|"),
		[]byte("complex128|"),
	},
	singleLineCommentStart: []byte("//"),
	multiLineCommentStart:  nil,
	multiLineCommentEnd:    nil,
	flags:                  HL_HIGHLIGHT_NUMBERS | HL_HIGHLIGHT_STRINGS,
}

var cSyntax = Syntax{
	fileType: []byte("C"),
	fileMatch: [][]byte{
		[]byte(".c"),
		[]byte(".h"),
		[]byte(".cpp"),
	},
	keywords: [][]byte{
		[]byte("switch"),
		[]byte("if"),
		[]byte("while"),
		[]byte("for"),
		[]byte("break"),
		[]byte("continue"),
		[]byte("return"),
		[]byte("else"),
		[]byte("struct"),
		[]byte("union"),
		[]byte("typedef"),
		[]byte("static"),
		[]byte("enum"),
		[]byte("class"),
		[]byte("case"),
		[]byte("int|"),
		[]byte("long|"),
		[]byte("double|"),
		[]byte("float|"),
		[]byte("char|"),
		[]byte("unsinged|"),
		[]byte("signed|"),
		[]byte("void|"),
	},
	singleLineCommentStart: []byte("//"),
	multiLineCommentStart:  []byte("/*"),
	multiLineCommentEnd:    []byte("*/"),
	flags:                  HL_HIGHLIGHT_NUMBERS | HL_HIGHLIGHT_STRINGS,
}

var syntaxDB = SyntaxDB{
	syntax: []Syntax{
		cSyntax,
		goSyntax,
	},
}

func selectSyntax() {
	editor.syntax = nil
	if editor.filename == "" {
		return
	}

	ext := "." + strings.Split(editor.filename, ".")[1]

	for _, s := range syntaxDB.syntax {
		for _, match := range s.fileMatch {
			if string(match) == ext {
				editor.syntax = &s

				for _, row := range editor.row {
					updateSyntax(&row)
				}

				return
			}
		}
	}
}

func isSperator(char int) bool {
	if unicode.IsSpace(rune(char)) || strings.ContainsRune(",.()+-/*=~%<>[];", rune(char)) {
		return true
	} else {
		return false
	}
}

func updateSyntax(row *EditorRow) {
	row.highlight = make([]byte, len(row.render))
	for i := range row.highlight {
		row.highlight[i] = HL_NORMAL
	}

	if editor.syntax == nil {
		return
	}

	keywords := editor.syntax.keywords

	scs := editor.syntax.singleLineCommentStart
	mcs := editor.syntax.multiLineCommentStart
	mce := editor.syntax.multiLineCommentEnd

	var scsLen int
	var mcsLen int
	var mceLen int

	if string(scs) != "" {
		scsLen = len(scs)
	} else {
		scsLen = 0
	}

	if string(mcs) != "" {
		mcsLen = len(mcs)
	} else {
		mcsLen = 0
	}

	if string(mce) != "" {
		mceLen = len(mce)
	} else {
		mceLen = 0
	}

	prevSep := true
	inString := 0
	inComment := 0

	if row.idx > 0 && editor.row[row.idx-1].hlOpenComment > 0 {
		inComment = 1
	}

	i := 0
	for i < len(row.render) {
		char := row.render[i]
		prevHl := HL_NORMAL // TODO check if line start

		// single line comments
		if scsLen > 0 && inString == 0 && inComment != 0 {
			trimmedString := strings.Trim(string(row.chars), " ")
			if strings.HasPrefix(trimmedString, string(scs)) {
				for i := range row.highlight {
					row.highlight[i] = HL_COMMENT
				}
				break
			}
		}

		// multi line comments
		if mcsLen != 0 && mceLen != 0 && inString == 0 {
			if inComment != 0 {
				row.highlight[i] = HL_MULTILINE_COMMENT
				if i+mceLen <= len(row.render) && string(row.render[i:i+mceLen]) == string(mce) {
					for k := range mceLen {
						row.highlight[i+k] = HL_MULTILINE_COMMENT
					}
					i += mceLen
				} else {
					i++
					continue
				}
			} else if i+mcsLen <= len(row.render) && string(row.render[i:i+mcsLen]) == string(mcs) {
				for k := range mcsLen {
					row.highlight[i+k] = HL_MULTILINE_COMMENT
				}
				i += mcsLen
				inComment = 1
				continue
			}
		}

		// strings
		if editor.syntax.flags&HL_HIGHLIGHT_STRINGS != 0 {
			if inString != 0 {
				row.highlight[i] = HL_STRING

				// dont end string with escaped '
				if char == '\\' && i+1 < len(row.render) {
					row.highlight[i+1] = HL_STRING
					i += 2
					continue
				}

				if int(char) == inString {
					inString = 0
				}

				prevSep = true
				i++
				continue
			} else {
				if char == '"' || char == '\'' {
					inString = int(char)
					row.highlight[i] = HL_STRING
					i++
					continue
				}
			}
		}

		if i > 0 {
			prevHl = int(row.highlight[i-1])
		}

		// numbers
		if editor.syntax.flags&HL_HIGHLIGHT_NUMBERS != 0 {
			if unicode.IsDigit(rune(char)) &&
				(prevSep || prevHl == HL_NUMBER) ||
				(char == '.' && prevHl == HL_NUMBER) {

				row.highlight[i] = HL_NUMBER
				i++
				prevSep = false
				continue

			}
		}

		// keywords
		if prevSep {
			for _, kw := range keywords {
				klen := len(kw)
				kw2 := klen > 0 && kw[klen-1] == '|'
				if kw2 {
					klen--
				}

				if i+klen <= len(row.render) && string(row.render[i:i+klen]) == string(kw[:klen]) && (i+klen == len(row.render) || isSperator(int(row.render[i+klen]))) {
					for k := range klen {
						row.highlight[i+k] = byte(func() int {
							if kw2 {
								return HL_KEYWORD2
							} else {
								return HL_KEYWORD1
							}
						}())
					}
					i += klen
					break
				}

				if kw != nil {
					prevSep = false
					continue
				}

			}
		}

		prevSep = isSperator(int(char))
		i++
	}

	changed := 0
	if row.hlOpenComment != inComment {
		changed = 1
	}

	if changed > 0 && row.idx+1 < editor.rows {
		updateSyntax(&editor.row[row.idx+1])
	}
}

func syntaxToColor(hl int) int {
	switch hl {
	case HL_NUMBER:
		return RED
	case HL_KEYWORD1:
		return GREEN
	case HL_KEYWORD2:
		return BLUE
	case HL_COMMENT, HL_MULTILINE_COMMENT:
		return MAGENTA
	case HL_STRING:
		return CYAN
	case HL_SEARCH_MATCH:
		return RED
	default:
		return WHITE
	}
}
