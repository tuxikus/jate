package editor

import (
	"io"
	"os"
	"syscall"
	"time"
)

const (
	// C- keys ignore case
	KEY_C_AT = 0 + iota
	KEY_C_A
	KEY_C_B
	KEY_C_C
	KEY_C_D
	KEY_C_E
	KEY_C_F
	KEY_C_G
	KEY_C_H
	KEY_C_I
	KEY_C_J
	KEY_C_K
	KEY_C_L
	KEY_C_M
	KEY_C_N
	KEY_C_O
	KEY_C_P
	KEY_C_Q
	KEY_C_R
	KEY_C_S
	KEY_C_T
	KEY_C_U
	KEY_C_V
	KEY_C_W
	KEY_C_X
	KEY_C_Y
	KEY_C_Z
	KEY_C_OB // opening bracket: [ = esc
	KEY_C_SLASH
	KEY_C_CB // closing bracket

	KEY_BACKSPACE = 127

	KEY_ARROW_LEFT = iota + 1000
	KEY_ARROW_RIGHT
	KEY_ARROW_UP
	KEY_ARROW_DOWN
	KEY_PAGE_UP
	KEY_PAGE_DOWN
	KEY_HOME
	KEY_END
	KEY_DELETE

	KEY_M_UPPER_A = iota + 2000
	KEY_M_UPPER_B
	KEY_M_UPPER_C
	KEY_M_UPPER_D
	KEY_M_UPPER_E
	KEY_M_UPPER_F
	KEY_M_UPPER_G
	KEY_M_UPPER_H
	KEY_M_UPPER_I
	KEY_M_UPPER_J
	KEY_M_UPPER_K
	KEY_M_UPPER_L
	KEY_M_UPPER_M
	KEY_M_UPPER_N
	KEY_M_UPPER_O
	KEY_M_UPPER_P
	KEY_M_UPPER_Q
	KEY_M_UPPER_R
	KEY_M_UPPER_S
	KEY_M_UPPER_T
	KEY_M_UPPER_U
	KEY_M_UPPER_V
	KEY_M_UPPER_W
	KEY_M_UPPER_X
	KEY_M_UPPER_Y
	KEY_M_UPPER_Z
	KEY_M_LOWER_A
	KEY_M_LOWER_B
	KEY_M_LOWER_C
	KEY_M_LOWER_D
	KEY_M_LOWER_E
	KEY_M_LOWER_F
	KEY_M_LOWER_G
	KEY_M_LOWER_H
	KEY_M_LOWER_I
	KEY_M_LOWER_J
	KEY_M_LOWER_K
	KEY_M_LOWER_L
	KEY_M_LOWER_M
	KEY_M_LOWER_N
	KEY_M_LOWER_O
	KEY_M_LOWER_P
	KEY_M_LOWER_Q
	KEY_M_LOWER_R
	KEY_M_LOWER_S
	KEY_M_LOWER_T
	KEY_M_LOWER_U
	KEY_M_LOWER_V
	KEY_M_LOWER_W
	KEY_M_LOWER_X
	KEY_M_LOWER_Y
	KEY_M_LOWER_Z

	KEY_M_COLON
)

const (
	KEY_BINIDNG_MODE_NORMAL = 0 + iota
	KEY_BINDING_MODE_EMACS
	KEY_BINDING_MODE_VI
)

const (
	VI_MODE_NORMAL = 0 + iota
	VI_MODE_INSERT
	VI_MODE_VISUAL
)

type KeyBindingMode int
type ViMode int

// in go chars are runes, so just integer (int32) values
func readKey() int {
	buf := make([]byte, 1)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				normalExit()
				// eagain = no data available right now, try again later
			} else if err == syscall.EAGAIN {
				continue
			} else {
				panicExit("readKey\n" + err.Error())
			}
		}
		// successfully read one byte
		if n == 1 {
			break
		}
	}

	c := buf[0]

	// if special key
	if c == '\x1b' {
		// set non block true for reading meta (alt) keys
		fd := int(os.Stdin.Fd())
		syscall.SetNonblock(fd, true)
		defer syscall.SetNonblock(fd, false)
		time.Sleep(1 * time.Millisecond)

		buf = make([]byte, 1)
		if n, err := os.Stdin.Read(buf); err != nil || n != 1 {
			return '\x1b'
		}
		seq0 := buf[0]

		// if next byte is [
		if seq0 == '[' {
			buf = make([]byte, 1)
			if n, err := os.Stdin.Read(buf); err != nil || n != 1 {
				return '\x1b'
			}
			seq1 := buf[0]

			// detect special keys:
			// page up:   \x1b[5~ => c = '\x1b'; seq0 = '['; seq1 = '5'; seq2 = '~'
			// page down: \x1b[5~ => c = '\x1b'; seq0 = '['; seq1 = '6'; seq2 = '~'
			if seq1 >= '0' && seq1 <= '9' {
				buf = make([]byte, 1)
				if n, err := os.Stdin.Read(buf); err != nil || n != 1 {
					return '\x1b'
				}
				seq2 := buf[0]

				if seq2 == '~' {
					switch seq1 {
					case '1':
						return KEY_HOME
					case '3':
						return KEY_DELETE
					case '4':
						return KEY_END
					case '5':
						return KEY_PAGE_UP
					case '6':
						return KEY_PAGE_DOWN
					case '7':
						return KEY_HOME
					case '8':
						return KEY_END
					}
				}
			} else {
				switch seq1 {
				case 'A':
					return KEY_ARROW_UP
				case 'B':
					return KEY_ARROW_DOWN
				case 'C':
					return KEY_ARROW_RIGHT
				case 'D':
					return KEY_ARROW_LEFT
				case 'H':
					return KEY_HOME
				case 'F':
					return KEY_END
				}
			}
		} else if seq0 == 'O' {
			buf = make([]byte, 1)
			if n, err := os.Stdin.Read(buf); err != nil || n != 1 {
				return '\x1b'
			}
			seq1 := buf[0]

			switch seq1 {
			case 'H':
				return KEY_HOME
			case 'F':
				return KEY_END
			}
		} else {
			// Meta keys
			// just esc with a letter
			// a -> \x1ba
			// b -> \x1bb
			// c -> \x1bc
			switch seq0 {
			case 'A':
				return KEY_M_UPPER_A
			case 'B':
				return KEY_M_UPPER_B
			case 'C':
				return KEY_M_UPPER_C
			case 'D':
				return KEY_M_UPPER_D
			case 'E':
				return KEY_M_UPPER_E
			case 'F':
				return KEY_M_UPPER_F
			case 'G':
				return KEY_M_UPPER_G
			case 'H':
				return KEY_M_UPPER_H
			case 'I':
				return KEY_M_UPPER_I
			case 'J':
				return KEY_M_UPPER_J
			case 'K':
				return KEY_M_UPPER_K
			case 'L':
				return KEY_M_UPPER_L
			case 'M':
				return KEY_M_UPPER_M
			case 'N':
				return KEY_M_UPPER_N
			case 'O':
				return KEY_M_UPPER_O
			case 'P':
				return KEY_M_UPPER_P
			case 'Q':
				return KEY_M_UPPER_Q
			case 'R':
				return KEY_M_UPPER_R
			case 'S':
				return KEY_M_UPPER_S
			case 'T':
				return KEY_M_UPPER_T
			case 'U':
				return KEY_M_UPPER_U
			case 'V':
				return KEY_M_UPPER_V
			case 'W':
				return KEY_M_UPPER_W
			case 'X':
				return KEY_M_UPPER_X
			case 'Y':
				return KEY_M_UPPER_Y
			case 'Z':
				return KEY_M_UPPER_Z
			// lower
			case 'a':
				return KEY_M_LOWER_A
			case 'b':
				return KEY_M_LOWER_B
			case 'c':
				return KEY_M_LOWER_C
			case 'd':
				return KEY_M_LOWER_D
			case 'e':
				return KEY_M_LOWER_E
			case 'f':
				return KEY_M_LOWER_F
			case 'g':
				return KEY_M_LOWER_G
			case 'h':
				return KEY_M_LOWER_H
			case 'i':
				return KEY_M_LOWER_I
			case 'j':
				return KEY_M_LOWER_J
			case 'k':
				return KEY_M_LOWER_K
			case 'l':
				return KEY_M_LOWER_L
			case 'm':
				return KEY_M_LOWER_M
			case 'n':
				return KEY_M_LOWER_N
			case 'o':
				return KEY_M_LOWER_O
			case 'p':
				return KEY_M_LOWER_P
			case 'q':
				return KEY_M_LOWER_Q
			case 'r':
				return KEY_M_LOWER_R
			case 's':
				return KEY_M_LOWER_S
			case 't':
				return KEY_M_LOWER_T
			case 'u':
				return KEY_M_LOWER_U
			case 'v':
				return KEY_M_LOWER_V
			case 'w':
				return KEY_M_LOWER_W
			case 'x':
				return KEY_M_LOWER_X
			case 'y':
				return KEY_M_LOWER_Y
			case 'z':
				return KEY_M_LOWER_Z

			case ':':
				return KEY_M_COLON
			}

		}
		// fallback
		return '\x1b'
	}
	// return a non escape character
	return int(c)
}

func processKeypress() {
	c := readKey()

	switch editor.keyBindingMode {
	case KEY_BINIDNG_MODE_NORMAL:
		processKeyPressNormal(c)
	case KEY_BINDING_MODE_EMACS:
		processKeyPressEmacs(c)
	case KEY_BINDING_MODE_VI:
		processKeyPressVi(c)
	}
}

func processKeyPressNormal(c int) {

}

func processKeyPressEmacs(c int) {
	switch c {
	case '\r':
		insertNewLine()

	case KEY_M_LOWER_M:
		moveCursorToIndentation()

	case KEY_M_LOWER_F:
		moveCursorWordForward()

	case KEY_M_LOWER_B:
		moveCursorWordBackward()

	case KEY_C_K:
		rowDeleteContent(editor.cursorY)

	case KEY_C_X:
		fileSave()

	case KEY_C_S:
		search()

	case KEY_M_COLON:
		executeCommand()

	case KEY_C_Q:
		normalExit()

	case KEY_BACKSPACE:
		deleteChar()

	case KEY_C_D, KEY_DELETE:
		moveCursorRight()
		deleteChar()

	case KEY_PAGE_UP:
		editor.cursorY = editor.rowOffset

		times := editor.screenRows
		for times > 0 {
			moveCursorUp()
			times--
		}

	case KEY_PAGE_DOWN:
		editor.cursorY = editor.rowOffset + editor.screenRows - 1
		if editor.cursorY > editor.rows {
			editor.cursorY = editor.rows
		}

		times := 0
		for times < editor.screenRows {
			moveCursorDown()
			times++
		}

	case KEY_C_E, KEY_END:
		if editor.cursorY < editor.rows {
			editor.cursorX = editor.row[editor.cursorY].length
		}

	case KEY_C_A, KEY_HOME:
		editor.cursorX = 0

	case KEY_C_N, KEY_ARROW_DOWN:
		moveCursorDown()

	case KEY_C_P, KEY_ARROW_UP:
		moveCursorUp()

	case KEY_C_B, KEY_ARROW_LEFT:
		moveCursorLeftEmacs()

	case KEY_C_F, KEY_ARROW_RIGHT:
		moveCursorRight()

	default:
		insertChar(c)
	}
}

///////////////////////////////////////////////////////////////////////////////
//                                  vi stuff                                 //
///////////////////////////////////////////////////////////////////////////////

func processKeyPressVi(c int) {
	switch editor.viMode {
	case VI_MODE_NORMAL:
		switch c {
		case 'h':
			moveCursorLeftVi()
		case 'l':
			moveCursorRightVi()
		case 'j':
			moveCursorDownVi()
		case 'k':
			moveCursorUpVi()
		case 'i':
			viEnableInsertMode()
			return
		case ':':
			executeCommand()
		}
	case VI_MODE_INSERT:
		switch c {
		case KEY_BACKSPACE:
			deleteChar()
		case '\x1b':
			viEnableNormalMode()
		default:
			insertChar(c)
		}
	case VI_MODE_VISUAL:
		switch c {
		case '\x1b':
			viEnableNormalMode()
		}

	}

}

func viEnableInsertMode() {
	editor.viMode = VI_MODE_INSERT
}

func viEnableNormalMode() {
	editor.viMode = VI_MODE_NORMAL
}

// func viEnableVisualMode() {
//	editor.viMode = VI_MODE_VISUAL
// }
