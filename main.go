package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term" // used to enable raw mode or get the terminal size, maybe change to syscalls directly
)

const VERSION = "0.0.1"
const TAB_WIDTH = 4
const EXIT_TRIES = 3

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

func insertRow(at int, s string) {
	if at < 0 || at > editor.rows {
		return
	}

	// new row
	editor.row = append(editor.row, EditorRow{})

	// shift rows
	copy(editor.row[at+1:], editor.row[at:])

	editor.row[at].chars = []byte(s)
	editor.row[at].length = len(editor.row[at].chars)
	editor.row[at].render = nil
	editor.row[at].renderLength = 0

	updateRow(&editor.row[at])

	editor.rows++
	editor.fileModified++
}

func updateRow(row *EditorRow) {
	// count tabs
	tabs := 0
	for _, char := range row.chars {
		if char == '\t' {
			tabs++
		}
	}

	size := len(row.chars) + tabs*(TAB_WIDTH-1) + 1
	row.render = make([]byte, 0, size)

	idx := 0
	for _, char := range row.chars {
		if char == '\t' {
			row.render = append(row.render, '#')
			idx++
			for idx%TAB_WIDTH != 0 {
				row.render = append(row.render, '#')
				idx++
			}
		} else {
			row.render = append(row.render, char)
			idx++
		}
	}

	row.renderLength = len(row.render)
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

// visual part of the editor
func drawRows(ab *AppendBuffer) {
	for y := range editor.screenRows {
		filerow := y + editor.rowOffset
		// print ~ after the file content
		if filerow >= editor.rows {
			// only display the welcome message if no file is loaded
			if editor.rows == 0 && y == editor.screenRows/2 {
				// message
				// draw tilde at start of line
				appendBufferAppend(ab, []byte("~"))
				// draw line 1
				messageLine1 := "jate - just another text editor"
				padding := ((editor.screenColumns - len(messageLine1)) / 2) - 1 // -1 = tilde
				for padding > 0 {
					appendBufferAppend(ab, []byte(" "))
					padding--
				}
				appendBufferAppend(ab, []byte(messageLine1))

				// fill line 1 so that line 2 is centered
				padding = (editor.screenColumns - len(messageLine1)) / 2
				for padding > 0 {
					appendBufferAppend(ab, []byte(" "))
					padding--
				}

				// draw tilde at start of line
				appendBufferAppend(ab, []byte("~"))
				// draw line 2
				messageLine2 := fmt.Sprintf("Version: %s", VERSION)
				padding = ((editor.screenColumns - len(messageLine2)) / 2) - 1 // -1 = tilde
				for padding > 0 {
					appendBufferAppend(ab, []byte(" "))
					padding--
				}
				appendBufferAppend(ab, []byte(messageLine2))
			} else {
				appendBufferAppend(ab, []byte("~"))
			}
		} else {
			max := editor.row[filerow].renderLength - editor.columnOffset
			appendBufferAppend(ab, editor.row[filerow].render[editor.columnOffset:max])
		}

		appendBufferAppend(ab, []byte("\x1b[K"))
		appendBufferAppend(ab, []byte("\r\n"))
	}
}

func deleteRow(at int) {
	if at < 0 || at >= editor.rows {
		return
	}

	// copy(dst, src)
	copy(editor.row[at:], editor.row[at+1:])
	editor.rows--
	editor.fileModified++
}

func rowAppendString(row *EditorRow, s string) {
	//row.chars = make([]byte, len(row.chars) + len(s))
	//copy(row.chars[:len(row.chars)], s)

	row.chars = append(row.chars, s...)
	row.length = len(row.chars)
	updateRow(row)
	editor.fileModified++
}

func rowDeleteChar(row *EditorRow, at int) {
	if at < 0 || at > row.length {
		return
	}

	// copy(dst, src)
	copy(row.chars[at:], row.chars[at+1:])
	row.chars = row.chars[:len(row.chars)-1]
	row.length--
	updateRow(row)
	editor.fileModified++
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

func rowInsertChar(row *EditorRow, at int, char byte) {
	if at < 0 || at > row.length {
		at = row.length
	}

	row.chars = append(row.chars, 0)       // add one char to make room for new char
	copy(row.chars[at+1:], row.chars[at:]) // shift all chars from at to the right
	row.chars[at] = char                   // add the char

	row.length++
	updateRow(row)
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

func rowsToString() string {
	s := ""
	for _, row := range editor.row {
		s += string(row.chars)
		s += "\n"
	}
	return s
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

// get the dimensions of the used terminal
func getTerminalSize() {
	columns, rows, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
	}
	editor.screenColumns = columns
	editor.screenRows = rows
}

var lastMatch = -1
var searchDirection = 1

func searchCallback(query []byte, key int) {
	if key == '\r' || key == '\x1b' {
		lastMatch = -1
		searchDirection = 1
		return
	} else if key == KEY_ARROW_RIGHT || key == KEY_ARROW_DOWN {
		searchDirection = 1
	} else if key == KEY_ARROW_RIGHT || key == KEY_ARROW_UP {
		searchDirection = -1
	} else {
		lastMatch = -1
		searchDirection = 1
	}

	if lastMatch == -1 {
		searchDirection = 1
	}

	current := lastMatch

	for range editor.row {
		current += searchDirection

		if current == -1 {
			current = editor.rows - 1
		} else if current == editor.rows {
			current = 0
		}

		row := &editor.row[current]

		s := string(row.chars)
		match := strings.Index(s, string(query))
		if match != -1 {
			lastMatch = current
			editor.cursorY = current
			editor.cursorX = renderXtoCursorX(row, match)
			editor.rowOffset = editor.rows
			break
		}
	}
}

func search() {
	oldCursorX := editor.cursorX
	oldCursorY := editor.cursorY
	oldColumnOffset := editor.columnOffset
	oldRowOffset := editor.rowOffset

	if prompt("Search: ", searchCallback) == nil {
		editor.cursorX = oldCursorX
		editor.cursorY = oldCursorY
		editor.columnOffset = oldColumnOffset
		editor.rowOffset = oldRowOffset
	}
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

func drawStatusBar(ab *AppendBuffer) {
	appendBufferAppend(ab, []byte("\x1b[7m"))

	left := fmt.Sprintf("File: %s", editor.filename)
	if editor.fileModified != 0 {
		left += " -modified-"
	}

	right := fmt.Sprintf("Lines: %d", editor.rows)

	appendBufferAppend(ab, []byte(left))
	for range editor.screenColumns - len(left) - len(right) {
		appendBufferAppend(ab, []byte(" "))
	}
	appendBufferAppend(ab, []byte(right))
	appendBufferAppend(ab, []byte("\x1b[m"))
	appendBufferAppend(ab, []byte("\r\n"))
}

func drawMessageBar(ab *AppendBuffer) {
	appendBufferAppend(ab, []byte("\x1b[K"))
	appendBufferAppend(ab, []byte(editor.statusMessage))
}

func setStatusMessage(format string, a ...interface{}) {
	editor.statusMessage = fmt.Sprintf(format, a...)
}

func initialize() {
	editor.cursorX = 0
	editor.cursorY = 0
	editor.renderX = 0
	editor.rows = 0
	editor.rowOffset = 0
	editor.columnOffset = 0
	editor.row = nil
	editor.filename = ""
	editor.statusMessage = ""
	editor.fileModified = 0

	getTerminalSize()

	editor.screenRows -= 2 // space for statusbar and status message
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
