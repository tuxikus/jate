package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

const (
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
	KEY_C_CB       // closing bracket
	KEY_BACKSPACE  = 127
	KEY_ARROW_LEFT = iota + 1000
	KEY_ARROW_RIGHT
	KEY_ARROW_UP
	KEY_ARROW_DOWN
	KEY_PAGE_UP
	KEY_PAGE_DOWN
	KEY_HOME
	KEY_END
	KEY_DELETE
)

var exitTries = 0

func processKeypress() {
	c := readKey()

	switch c {
	case '\r':
		insertNewLine()

	case KEY_C_X:
		fileSave()

	case KEY_C_S:
		search()

	case KEY_C_Q:
		if editor.fileModified != 0 && exitTries < EXIT_TRIES {
			setStatusMessage(fmt.Sprintf("File modified, exit without saving? Press C-q %d more times", EXIT_TRIES-exitTries))
			exitTries++
			return
		}
		normalExit()

	case KEY_BACKSPACE:
		deleteChar()

	case KEY_DELETE:
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
		moveCursorLeft()

	case KEY_C_F, KEY_ARROW_RIGHT:
		moveCursorRight()

	default:
		insertChar(c)
	}

	exitTries = 0
}

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
				panicExit("readKey")
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
		buf = make([]byte, 1)
		if n, err := os.Stdin.Read(buf); err != nil || n != 1 {
			return '\x1b'
		}
		seq0 := buf[0]

		buf = make([]byte, 1)
		if n, err := os.Stdin.Read(buf); err != nil || n != 1 {
			return '\x1b'
		}
		seq1 := buf[0]

		// if next byte is [
		if seq0 == '[' {
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
			switch seq1 {
			case 'H':
				return KEY_HOME
			case 'F':
				return KEY_END
			}
		}
		// fallback
		return '\x1b'
	}
	// return a non escape character
	return int(c)
}
