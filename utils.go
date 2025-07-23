package main

import (
	"fmt"
	"os"
)

func normalExit() {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H"))  // move cursor to 1 1
	disableRawMode()
	printEditorStuff()
	os.Exit(0)
}

func panicExit(message string) {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H"))  // move cursor to 1 1
	disableRawMode()
	fmt.Println(message)
	os.Exit(1)
}
