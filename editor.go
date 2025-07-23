package main

import (
	"golang.org/x/term" // used to enable raw mode or get the terminal size, maybe change to syscalls directly
)

// type to store information about a row (line)
// chars  + length is the content
// render + renderLength is the rendered content
type EditorRow struct {
	length       int
	renderLength int
	chars        []byte
	render       []byte
}

// type to store global editor stuff
type Editor struct {
	cursorX       int
	cursorY       int
	renderX       int
	rowOffset     int // index of row[]
	columnOffset  int // index of row.chars[]
	screenRows    int
	screenColumns int
	rows          int
	row           []EditorRow
	filename      string
	fileModified  int
	statusMessage string
	oldTermState  *term.State // used to restore the terminal config after enabling raw mode
}

var editor = Editor{
	screenRows:    0,
	screenColumns: 0,
	oldTermState:  nil,
}
