package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

const (
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

// used for simple debugging
func printEditorStuff() {
	fmt.Println("editor.filename =", editor.filename)
	fmt.Println("editor.cursorX =", editor.cursorX)
	fmt.Println("editor.cursorY =", editor.cursorY)
	fmt.Println("editor.rows =", editor.screenRows)
	fmt.Println("editor.columns =", editor.screenColumns)
	fmt.Println("editor.rows =", editor.rows)
	fmt.Println("editor.row.chars =")
	for _, line := range editor.row {
		fmt.Println(line.chars, line.length)
	}
	fmt.Println("editor.row.render =")
	for _, line := range editor.row {
		fmt.Println(line.render, line.renderLength)
	}
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

func scroll() {
	editor.renderX = 0
	if editor.cursorY < editor.rows {
		editor.renderX = cursorXToRenderX(&editor.row[editor.cursorY], editor.cursorX)
	}

	if editor.cursorY < editor.rowOffset {
		editor.rowOffset = editor.cursorY
	}

	if editor.cursorY >= editor.rowOffset+editor.screenRows {
		editor.rowOffset = editor.cursorY - editor.screenRows + 1
	}

	if editor.renderX < editor.columnOffset {
		editor.columnOffset = editor.renderX
	}

	if editor.renderX >= editor.columnOffset+editor.screenColumns {
		editor.columnOffset = editor.renderX - editor.screenColumns + 1
	}
}

func deleteChar() {
	// last line + 1
	if editor.cursorY == editor.rows {
		return
	}

	// starting position
	if editor.cursorX == 0 && editor.cursorY == 0 {
		return
	}

	row := &editor.row[editor.cursorY]
	if editor.cursorX > 0 {
		rowDeleteChar(&editor.row[editor.cursorY], editor.cursorX-1)
		editor.cursorX--
	} else {
		// cursor on the beginning of the line => delet this line and append to line above
		editor.cursorX = editor.row[editor.cursorY-1].length
		rowAppendString(&editor.row[editor.cursorY-1], string(row.chars))
		deleteRow(editor.cursorY)
		editor.cursorY--
	}
}

func insertChar(c int) {
	if editor.cursorY == editor.rows {
		insertRow(editor.rows, "")
	}
	rowInsertChar(&editor.row[editor.cursorY], editor.cursorX, byte(c))
	editor.cursorX++
	editor.fileModified++
}

func moveCursor(key int) {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}

	switch key {
	case KEY_ARROW_LEFT:
		if editor.cursorX != 0 {
			editor.cursorX--
		} else if editor.cursorY > 0 {
			editor.cursorY--
			editor.cursorX = editor.row[editor.cursorY].length
		}
	case KEY_ARROW_RIGHT:
		if row != nil && editor.cursorX < row.length {
			editor.cursorX++
		} else if row != nil && editor.cursorX == row.length {
			editor.cursorY++
			editor.cursorX = 0
		}
	case KEY_ARROW_UP:
		if editor.cursorY > 0 {
			editor.cursorY--
		}
	case KEY_ARROW_DOWN:
		if editor.cursorY < editor.rows {
			editor.cursorY++
		}
	}

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		// get the new row if y changed
		row = &editor.row[editor.cursorY]
	}

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
		}
	}
}

func renderXtoCursorX(row *EditorRow, renderX int) int {
	currentRenderX := 0

	for cursorX := 0; cursorX < len(row.chars); cursorX++ {
		if row.chars[cursorX] == '\t' {
			currentRenderX += (TAB_WIDTH - 1) - (currentRenderX % TAB_WIDTH)
		}
		currentRenderX++

		if currentRenderX > renderX {
			return cursorX
		}
	}

	return len(row.chars)
}

func cursorXToRenderX(row *EditorRow, cursorX int) int {
	renderX := 0
	for i := range cursorX {
		if row.chars[i] == '\t' {
			// how many columns right to the last tab
			renderX += TAB_WIDTH - 1 - (renderX % TAB_WIDTH)
		}
		renderX++
	}

	return renderX
}

func refreshScreen() {
	var appendBuffer AppendBuffer

	scroll()

	// hide the cursor
	appendBufferAppend(&appendBuffer, []byte("\x1b?25l"))

	// clear the screen
	appendBufferAppend(&appendBuffer, []byte("\x1b[2J"))

	// reposition the cursor to the beginning
	// H: VT100 cursor position
	// [10;10H move cursor to row 10 and column 10
	// default is 1;1
	appendBufferAppend(&appendBuffer, []byte("\x1b[H"))

	drawRows(&appendBuffer)       // screenrows - 2
	drawStatusBar(&appendBuffer)  // screenrows - 1
	drawMessageBar(&appendBuffer) // screenrows

	cursorVt100 := fmt.Sprintf("\x1b[%d;%dH", editor.cursorY-editor.rowOffset+1, editor.renderX-editor.columnOffset+1)
	appendBufferAppend(&appendBuffer, []byte(cursorVt100))

	// show the cursor
	appendBufferAppend(&appendBuffer, []byte("\x1b[?25h"))

	os.Stdin.Write(appendBuffer.chars) // the only write call per refresh
}

var exitTries = 0

func processKeypress() {
	c := readKey()

	switch c {
	case '\r':
		insertNewLine()

	// C-s
	case 19:
		fileSave()

	// C-f
	case 6:
		search()

	// C-q
	case 17:
		if editor.fileModified != 0 && exitTries < EXIT_TRIES {
			setStatusMessage(fmt.Sprintf("File modified, exit without saving? Press C-q %d more times", EXIT_TRIES-exitTries))
			exitTries++
			return
		}
		normalExit()

	// TODO add C-h
	case KEY_BACKSPACE:
		deleteChar()

	case KEY_DELETE:
		moveCursor(KEY_ARROW_RIGHT)
		deleteChar()

	case KEY_PAGE_UP:
		editor.cursorY = editor.rowOffset

		times := editor.screenRows
		for times > 0 {
			moveCursor(KEY_ARROW_UP)
			times--
		}
	case KEY_PAGE_DOWN:
		editor.cursorY = editor.rowOffset + editor.screenRows - 1
		if editor.cursorY > editor.rows {
			editor.cursorY = editor.rows
		}

		times := 0
		for times < editor.screenRows {
			moveCursor(KEY_ARROW_DOWN)
			times++
		}
	case KEY_END:
		if editor.cursorY < editor.rows {
			editor.cursorX = editor.row[editor.cursorY].length
		}
	case KEY_HOME:
		editor.cursorX = 0

	case KEY_ARROW_DOWN, KEY_ARROW_LEFT, KEY_ARROW_RIGHT, KEY_ARROW_UP:
		moveCursor(c)

	// C-l
	case 12:
		break
	default:
		insertChar(c)
	}

	exitTries = 0
}

func insertNewLine() {
	if editor.cursorX == 0 {
		insertRow(editor.cursorY, "")
	} else {
		row := &editor.row[editor.cursorY]
		// insert new row
		insertRow(editor.cursorY+1, string(row.chars[editor.cursorX:]))

		// edit old line
		row = &editor.row[editor.cursorY]
		// length is now the line break point
		row.length = editor.cursorX
		// chars are all up to the cursor location
		row.chars = row.chars[:editor.cursorX]
		updateRow(row)
	}
	editor.cursorY++
	editor.cursorX = 0
}

func main() {
	enableRawMode()
	initialize()

	if len(os.Args) > 1 {
		fileOpen(os.Args[1])
	}

	setStatusMessage("C-q to quit")

	for {
		refreshScreen()
		processKeypress()
	}
}
