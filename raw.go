package main

import (
	"golang.org/x/term" // used to enable raw mode or get the terminal size, maybe change to syscalls directly
	"os"
)

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
