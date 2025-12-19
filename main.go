package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"time"
	"unicode"

	"golang.org/x/term" // used to enable raw mode or get the terminal size, maybe change to syscalls directly
)

const (
	Version  = "1.0.0"
	TabWidth = 4
)

// type to store information about a row (line)
// chars  + length is the content
// render + renderLength is the rendered content
type EditorRow struct {
	length       int
	renderLength int
	chars        []byte
	render       []byte
	idx          int
}

// type to store global editor stuff
type Editor struct {
	cursorX        int
	cursorY        int
	renderX        int
	rowOffset      int // index of row[]
	columnOffset   int // index of row.chars[]
	screenRows     int
	screenColumns  int
	rows           int
	row            []EditorRow
	keymapBindings keymapBindings
	filename       string
	fileModified   int
	statusMessage  string
	oldTermState   *term.State // used to restore the terminal config after enabling raw mode
}

type AppendBuffer struct {
	chars []byte
}

var editor = Editor{
	screenRows:    0,
	screenColumns: 0,
	oldTermState:  nil,
}

func (e *Editor) setStatusMessage(msg string) {
	e.statusMessage = msg
}

func runEditor() {
	for {
		draw()
		processKeypress()
	}
}

func initEditor() {
	enableRawMode()

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
	editor.keymapBindings = initNormalKeymapBindings()

	columns, rows := getTerminalSize()
	editor.screenColumns = columns
	editor.screenRows = rows

	editor.screenRows -= 2 // space for statusbar and status message
}

// get the dimensions of the used terminal
func getTerminalSize() (int, int) {
	columns, rows, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panicExit("getTerminalSize\n" + err.Error())
	}

	return columns, rows
}

func enableRawMode() {
	var err error
	editor.oldTermState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
}

func disableRawMode() {
	term.Restore(int(os.Stdin.Fd()), editor.oldTermState)
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
		rowDelete(editor.cursorY)
		editor.cursorY--
	}
}

func insertChar(c int) {
	if c < 32 || c > 126 {
		return
	}

	if editor.cursorY == editor.rows {
		insertRow(editor.rows, "")
	}
	rowInsertChar(&editor.row[editor.cursorY], editor.cursorX, byte(c))
	editor.cursorX++
	editor.fileModified++
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

//
// AppendBuffer
//

func appendBufferAppend(ab *AppendBuffer, chars []byte) {
	ab.chars = append(ab.chars, chars...)
}

//
// Cursor
//

func moveCursorUp() {
	if editor.cursorY > 0 {
		editor.cursorY--
	}

	currentRow := getCurrentRow()

	// check if cursor is past the row length
	if currentRow != nil {
		if editor.cursorX > currentRow.length {
			editor.cursorX = currentRow.length
		}
	}

}

func moveCursorDown() {
	if editor.cursorY < editor.rows {
		editor.cursorY++
	}

	currentRow := getCurrentRow()

	// check if cursor is past the row length
	if currentRow != nil {
		if editor.cursorX > currentRow.length {
			editor.cursorX = currentRow.length
		}
	}
}

func moveCursorLeft() {
	if editor.cursorX != 0 {
		editor.cursorX--
	} else if editor.cursorY > 0 {
		editor.cursorY--
		editor.cursorX = len(editor.row[editor.cursorY].chars)
	}

	currentRow := getCurrentRow()

	// check if cursor is past the row length
	if currentRow != nil {
		if editor.cursorX > currentRow.length {
			editor.cursorX = currentRow.length
		}
	}
}

func moveCursorRight() {
	currentRow := getCurrentRow()

	if currentRow != nil && editor.cursorX < currentRow.length {
		editor.cursorX++
	} else if currentRow != nil && editor.cursorX == currentRow.length {
		editor.cursorY++
		editor.cursorX = 0
	}

	currentRow = getCurrentRow()

	// check if cursor is past the row length
	if currentRow != nil {
		if editor.cursorX > currentRow.length {
			editor.cursorX = currentRow.length
		}
	}
}

func moveCursorToBeginning() {
	editor.cursorX = 0
}

func moveCursorToEnd() {
	editor.cursorX = len(getCurrentRow().chars)
}

//
// Draw
//

func draw() {
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
				messageLine2 := fmt.Sprintf("Version: %s", Version)
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
			drawStart := min(editor.columnOffset, len(editor.row[filerow].render))
			drawEnd := min(drawStart+editor.screenColumns, len(editor.row[filerow].render))

			appendBufferAppend(ab, editor.row[filerow].render[drawStart:drawEnd])
			appendBufferAppend(ab, []byte("\x1b[39m"))
		}

		appendBufferAppend(ab, []byte("\x1b[K"))
		appendBufferAppend(ab, []byte("\r\n"))
	}
}

//
// FileIO
//

func fileOpen(filename string) {
	if editor.filename == "" {
		editor.filename = filename
	} else {
		// if file already loaded, reset editor
		// TODO function
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

		editor.filename = filename
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		panicExit("open\n" + err.Error())
	}

	var contentAsBytes []byte

	for _, char := range content {
		if char == '\n' || char == 10 {
			insertRow(editor.rows, string(contentAsBytes))
			contentAsBytes = nil
			continue
		}
		contentAsBytes = append(contentAsBytes, char)
	}
	editor.fileModified = 0
}

//
//
//

func fileSave() {
	if editor.filename == "" {
		editor.filename = string(prompt("Save as: ", nil))
	}

	file, err := os.OpenFile(editor.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panicExit("save " + err.Error())
	}
	defer file.Close()

	fileString := rowsToString()
	fileBytes := []byte(fileString)
	file.Write(fileBytes)
	setStatusMessage("%d bytes saved to disk", len(fileBytes))
	editor.fileModified = 0
}

//
// key
//

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

//
// keymap
//

const (
	KEYMAP_NORMAL = iota
)

type keymapBindings map[int]func()

var normalKeymapBindings keymapBindings

func initNormalKeymapBindings() keymapBindings {
	return keymapBindings{
		13:              insertNewLine, // 13 => Carriage Return (Enter)
		KEY_ARROW_UP:    moveCursorUp,
		KEY_ARROW_RIGHT: moveCursorRight,
		KEY_ARROW_DOWN:  moveCursorDown,
		KEY_ARROW_LEFT:  moveCursorLeft,
		KEY_HOME:        moveCursorToBeginning,
		KEY_END:         moveCursorToEnd,
		KEY_BACKSPACE:   deleteChar,
		KEY_C_C:         normalExit,
		KEY_C_S:         fileSave,
		KEY_C_F:         search,
	}
}

//
// menu bar
//

func drawMessageBar(ab *AppendBuffer) {
	appendBufferAppend(ab, []byte("\x1b[K"))
	appendBufferAppend(ab, []byte(editor.statusMessage))
}

//
// prompt
//

type PromptCallback func(input []byte, key int)

func prompt(prompt string, promptCallback PromptCallback) []byte {
	buf := make([]byte, 0)

	for {
		setStatusMessage("%s%s", prompt, buf)
		draw()

		c := readKey()

		if c == KEY_BACKSPACE {
			if len(buf) != 0 {
				buf = buf[:len(buf)-1]
			}
		} else if c == '\x1b' {
			setStatusMessage("")
			if promptCallback != nil {
				promptCallback(buf, c)
			}
			return nil
		} else if c == '\r' || c == '\n' {
			if len(buf) != 0 {
				setStatusMessage("")
				return buf[:]
			}
			if promptCallback != nil {
				promptCallback(buf, c)
			}
		} else if !unicode.IsControl(rune(c)) && c < 128 {
			buf = append(buf, byte(c))
			// completion
		}

		if promptCallback != nil {
			promptCallback(buf, c)
		}
	}
}

//
// row
//

// TODO: write test
func getCurrentRow() *EditorRow {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}

	return row
}

// deletes a row from editor.row
func rowDelete(at int) {
	if at < 0 || at >= editor.rows {
		return
	}

	// copy(dst, src)
	// copy all rows below at to the index of at
	copy(editor.row[at:], editor.row[at+1:])

	for i := at; i < editor.rows-1; i++ {
		editor.row[i].idx--
	}

	editor.rows--
	editor.fileModified++
}

func rowAppendString(row *EditorRow, s string) {
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

func rowsToString() string {
	s := ""
	for i := range editor.rows {
		s += string(editor.row[i].chars) + "\n"
	}

	return s
}

func insertRow(at int, s string) {
	if at < 0 || at > editor.rows {
		return
	}

	// new empty row
	editor.row = append(editor.row, EditorRow{})

	// shift rows
	copy(editor.row[at+1:], editor.row[at:])

	for i := at + 1; i <= editor.rows; i++ {
		editor.row[i].idx++
	}

	editor.row[at].idx = at

	editor.row[at].chars = []byte(s)
	editor.row[at].length = len(editor.row[at].chars)
	editor.row[at].render = nil
	editor.row[at].renderLength = 0

	updateRow(&editor.row[at])

	editor.rows++
	editor.fileModified++
}

// build the row.render from row.chars
func updateRow(row *EditorRow) {
	// count tabs
	tabs := 0
	for _, char := range row.chars {
		if char == '\t' {
			tabs++
		}
	}

	// TAB_WIDTH - 1 -> /t already a char
	size := len(row.chars) + tabs*(TabWidth-1)
	row.render = make([]byte, 0, size)

	idx := 0
	for _, char := range row.chars {
		if char == '\t' {
			row.render = append(row.render, ' ')
			idx++
			// if char is a tab check idx and add needed spaces to fill the tab
			for idx%TabWidth != 0 {
				row.render = append(row.render, ' ')
				idx++
			}
		} else {
			row.render = append(row.render, char)
			idx++
		}
	}

	row.renderLength = len(row.render)
}

// status bar

// TODO: move logic
func drawStatusBar(ab *AppendBuffer) {
	appendBufferAppend(ab, []byte("\x1b[7m"))

	var fType []byte
	var fName []byte
	var viMode []byte

	fType = []byte("")

	if editor.filename != "" {
		fName = []byte(editor.filename)
	} else {
		fName = []byte("-")
	}

	left := fmt.Sprintf(" %s %s File: %s Lines: %d:%d", viMode, fType, fName, editor.rows, editor.cursorY+1)
	if editor.fileModified != 0 {
		left += " -modified-"
	}

	t := time.Now().Format("15:04")
	right := fmt.Sprintf("%s  ", t)

	appendBufferAppend(ab, []byte(left))
	for range editor.screenColumns - len(left) - len(right) {
		appendBufferAppend(ab, []byte(" "))
	}
	appendBufferAppend(ab, []byte(right))
	appendBufferAppend(ab, []byte("\x1b[m"))
	appendBufferAppend(ab, []byte("\r\n"))
}

//
// search
//

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

		switch current {
		case -1:
			current = editor.rows - 1
		case editor.rows:
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

//
// status msg
//

func setStatusMessage(format string, a ...interface{}) {
	editor.setStatusMessage(fmt.Sprintf(format, a...))
}

//
// utils
//

func normalExit() {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H"))  // move cursor to 1 1
	disableRawMode()

	os.Exit(0)
}

func panicExit(message string) {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H"))  // move cursor to 1 1
	disableRawMode()
	fmt.Println(message)
	os.Exit(1)
}

func renderXtoCursorX(row *EditorRow, renderX int) int {
	currentRenderX := 0

	for cursorX := 0; cursorX < len(row.chars); cursorX++ {
		if row.chars[cursorX] == '\t' {
			currentRenderX += (TabWidth - 1) - (currentRenderX % TabWidth)
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
			renderX += TabWidth - 1 - (renderX % TabWidth)
		}
		renderX++
	}

	return renderX
}

func main() {
	initEditor()

	if len(os.Args) > 1 {
		fileOpen(os.Args[1])
	}

	runEditor()
}
