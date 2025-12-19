package main

import (
	"os"
)

func main() {
	Initialize()

	if len(os.Args) > 1 {
		FileOpen(os.Args[1])
	}

	Run()
}
