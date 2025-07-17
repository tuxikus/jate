package main

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/term" // used to enable raw mode or get the terminal size, maybe change to syscalls directly
)

const VERSION = "0.0.0"

const (
		KEY_ARROW_LEFT = 1_000
		KEY_ARROW_RIGHT = 1_001
		KEY_ARROW_UP = 1_002
		KEY_ARROW_DOWN = 1_003
		KEY_PAGE_UP = 2_000
		KEY_PAGE_DOWN = 2_001
		KEY_HOME = 3_000
		KEY_END = 3_001
		KEY_DELETE = 4_000
)

// struct to hold global editor stuff
type Editor struct {
		CursorX int
		CursorY int
		Rows int
		Columns int
		OldTermState *term.State // used to restore the editor config after enabling raw mode
}

// used to call write only once per refresh
type AppendBuffer struct {
		chars []byte
}

var editor = Editor{
		Rows: 0,
		Columns: 0,
		OldTermState: nil,
}

func appendBufferAppend(ab *AppendBuffer, chars []byte) {
		ab.chars = append(ab.chars, chars...) // ... is unpacking
}

func printEditorStuff() {
		fmt.Println("editor.CursorX =", editor.CursorX)
		fmt.Println("editor.CursorY =", editor.CursorY)
		fmt.Println("editor.Rows =", editor.Rows)
		fmt.Println("editor.Columns =", editor.Columns)
}

func normalExit() {
		os.Stdin.Write([]byte("\x1b[2J")) // clear
		os.Stdin.Write([]byte("\x1b[H")) // move cursor to 1 1
		disableRawMode()
		printEditorStuff()
		os.Exit(0)
}

func panicExit(message string) {
		os.Stdin.Write([]byte("\x1b[2J")) // clear
		os.Stdin.Write([]byte("\x1b[H")) // move cursor to 1 1
		disableRawMode()
		fmt.Println(message)
		os.Exit(1)
}

func enableRawMode() {
		var err error
		editor.OldTermState, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
				panic(err)
		}
}

func disableRawMode() {
		term.Restore(int(os.Stdin.Fd()), editor.OldTermState)
}

// in go chars in '' are runes, so just integer (int32) values
func readKey() int {
		// use byte instead of byte slice with len 1
		// ???
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

		// if start with escape
		if c == '\x1b' {
				// check the next two bytes
				seq := make([]byte, 2)
				n1, err1 := os.Stdin.Read(seq[:1])
				if err1 != nil || n1 != 1 {
						return '\x1b'
				}
				n2, err2 := os.Stdin.Read(seq[1:])
				if err2 != nil || n2 != 1 {
						return '\x1b'
				}

				// if next byte is [
				if seq[0] == '[' {
						switch seq[1] {
						case 'A':
								return KEY_ARROW_UP
						case 'B':
								return KEY_ARROW_DOWN
						case 'C':
								return KEY_ARROW_RIGHT
						case 'D':
								return KEY_ARROW_LEFT
						}
				}
				// fallback
				return '\x1b'
		}

		// return a non escape character
		return int(c)
}

func drawRows(ab *AppendBuffer) {
		for y := range editor.Rows {
				if y == editor.Rows / 2 {
						message := fmt.Sprintf("Hello - Version: %s", VERSION)
						padding := (editor.Columns - len(message)) / 2
						for padding > 0 {
								appendBufferAppend(ab, []byte(" "))
								padding--
						}
						appendBufferAppend(ab, []byte(message))
				} else {
						appendBufferAppend(ab, []byte("~"))
				}

				appendBufferAppend(ab, []byte("\x1b[K"))
				if y < editor.Rows - 1 {
						appendBufferAppend(ab, []byte("\r\n"))
				}
		}
}

func moveCursor(key int) {
		switch key {
		case KEY_ARROW_LEFT:
				if editor.CursorX > 0 {
						editor.CursorX--
				}
		case KEY_ARROW_RIGHT:
				if editor.CursorX < editor.Columns {
						editor.CursorX++
				}
		case KEY_ARROW_UP:
				if editor.CursorY > 0 {
						editor.CursorY--
				}
		case KEY_ARROW_DOWN:
				if editor.CursorY < editor.Rows {
						editor.CursorY++
				}
		}
}

func refreshScreen() {
		var appendBuffer AppendBuffer

		// hide the cursor
		appendBufferAppend(&appendBuffer, []byte("\x1b?25l"))
		appendBufferAppend(&appendBuffer, []byte("\x1b[2J"))
		// clear the screen
		//appendBufferAppend(&appendBuffer, []byte("\x1b[2J"))
		// reposition the cursor to the beginning
		// H: VT100 cursor position
		// [10;10H move cursor to row 10 and column 10
		// default is 1;1
		appendBufferAppend(&appendBuffer, []byte("\x1b[H"))

		drawRows(&appendBuffer)

		cursorVt100 := fmt.Sprintf("\x1b[%d;%dH", editor.CursorY + 1, editor.CursorX + 1)
		appendBufferAppend(&appendBuffer, []byte(cursorVt100))

		// show the cursor
		appendBufferAppend(&appendBuffer, []byte("\x1b[?25h"))

		os.Stdin.Write(appendBuffer.chars) // the only write call per refresh
}

func processKey() {
		c := readKey()

		switch c {
		// C-q
		case 17:
				normalExit()

		case KEY_ARROW_DOWN, KEY_ARROW_LEFT, KEY_ARROW_RIGHT, KEY_ARROW_UP:
				moveCursor(c)

		default:
				fmt.Println(c)
		}
}

// get the dimensions of the used terminal
func getTerminalSize() {
		columns, rows, err := term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
				fmt.Println(err)
		}
		editor.Columns = columns
		editor.Rows = rows
}

func initialize() {
		editor.CursorX = 0
		editor.CursorY = 0
		getTerminalSize()
}

func main() {
		enableRawMode()
		initialize()


		for {
				refreshScreen()
				processKey()
		}
}
