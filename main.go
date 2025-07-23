package main

import (
	"os"
)

func main() {
	enableRawMode()
	initialize()

	if len(os.Args) > 1 {
		fileOpen(os.Args[1])
	}

	setStatusMessage("C-q to quit")

	for {
		draw()
		processKeypress()
	}
}
