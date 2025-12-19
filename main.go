package main

import (
	"os"
)

const (
	Version  = "0.0.1"
	TabWidth = 4
)

func main() {
	initEditor()

	if len(os.Args) > 1 {
		fileOpen(os.Args[1])
	}

	runEditor()
}
