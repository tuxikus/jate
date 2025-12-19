package main

import (
	"io"
	"os"
	"syscall"
	"time"
)

const (
	// C- keys ignore case
	KEY_C_AT = 0
	KEY_C_C  = 3
	KEY_C_F  = 6
	KEY_C_S  = 19

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
		// used for Meta + <key>. Unused
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
		switch seq0 {
		case '[':
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
		case 'O':
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
		default:
			return '\x1b'
		}
	}
	// return a non escape character
	return int(c)
}

func processKeypress() {
	c := readKey()

	if action, exists := editor.keymapBindings[c]; exists {
		if action != nil {
			action()
		}
		return
	} else {
		insertChar(c)
	}

}
