package editor

import (
	"os"

	"golang.org/x/term" // used to enable raw mode or get the terminal size, maybe change to syscalls directly
)

// type to store information about a row (line)
// chars  + length is the content
// render + renderLength is the rendered content
type EditorRow struct {
	length        int
	renderLength  int
	chars         []byte
	render        []byte
	highlight     []byte
	idx           int
	hlOpenComment int
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
	viMode         ViMode
	filename       string
	fileModified   int
	statusMessage  string
	syntax         *Syntax
	oldTermState   *term.State // used to restore the terminal config after enabling raw mode
}

var editor = Editor{
	screenRows:    0,
	screenColumns: 0,
	oldTermState:  nil,
}

func (e *Editor) setStatusMessage(msg string) {
	e.statusMessage = msg
}

func Run() {
	for {
		draw()
		processKeypress()
	}
}

func Initialize() {
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
	editor.syntax = nil
	editor.keymapBindings = emacsKeymapBindings

	columns, rows := getTerminalSize()
	editor.screenColumns = columns
	editor.screenRows = rows

	editor.screenRows -= 2 // space for statusbar and status message
}

// get the dimensions of the used terminal
func getTerminalSize() (int, int) {
	columns, rows, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		// TODO
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
