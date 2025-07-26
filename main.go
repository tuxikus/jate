package main

import (
	"jate/editor"
	"os"
)

func main() {
	editor.Initialize()

	if len(os.Args) > 1 {
		editor.FileOpen(os.Args[1])
	}

	editor.Run()
}
