package main

import (
	"fmt"
	"golang.org/x/term" // used to enable raw mode or get the terminal size, maybe change to syscalls directly
	"io"
	"os"
)

// struct to hold global editor stuff
type Editor struct {
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

func normalExit() {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H")) // move cursor to 1 1
	disableRawMode()
	fmt.Println(editor.Rows)
	os.Exit(0)
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

func readKey() byte {
	buf := make([]byte, 1)
	_, err := os.Stdin.Read(buf)
	if err == io.EOF {
		normalExit()
	}
	return buf[0]
}

func drawRows(ab *AppendBuffer) {
	for i := range editor.Columns {
		appendBufferAppend(ab, []byte("~"))
		if i < editor.Columns - 1 {
			appendBufferAppend(ab, []byte("\r\n"))
		}
	}
}

func refreshScreen() {
	var appendBuffer AppendBuffer

	// clear the screen
	appendBufferAppend(&appendBuffer, []byte("\x1b[2J"))
	// reposition the cursor to the beginning
	// H: VT100 cursor position
	// [10;10H move cursor to row 10 and column 10
	// default is 1;1
	appendBufferAppend(&appendBuffer, []byte("\x1b[H"))

	drawRows(&appendBuffer)

	appendBufferAppend(&appendBuffer, []byte("\x1b[H"))

	os.Stdin.Write(appendBuffer.chars) // the only write call per refresh
}

func processKey() {
	c := readKey()

	switch c {
	// C-q
	case 17:
		normalExit()
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
