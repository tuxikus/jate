package main

import (
	"fmt"
	"os"
)

func NormalExit() {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H"))  // move cursor to 1 1
	DisableRawMode()
	printEditorStuff()
	os.Exit(0)
}

func PanicExit(message string) {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H"))  // move cursor to 1 1
	DisableRawMode()
	fmt.Println(message)
	os.Exit(1)
}
