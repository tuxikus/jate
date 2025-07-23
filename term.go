package main

import (
	"golang.org/x/term" // used to enable raw mode or get the terminal size, maybe change to syscalls directly
	"os"
)

// get the dimensions of the used terminal
func getTerminalSize() {
	columns, rows, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		// TODO
	}
	editor.screenColumns = columns
	editor.screenRows = rows
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
