package main

import (
	"fmt"
	"golang.org/x/term" // used to enable raw mode, maybe change to syscalls directly
	"io"
	"os"
)

// struct to hold global editor stuff
type Editor struct {
	OldTermState *term.State // used to restore the editor config after enabling raw mode
}

var editor = Editor{
	OldTermState: nil,
}

func enableRawMode() {
	var err error
	editor.OldTermState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
}

func disableRawMode() {
	term.Restore(int(os.Stdin.Fd()), editor.OldTermState)
}

func main() {
	enableRawMode()
	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		if buf[0] == 113 {
			break
		}

		fmt.Println(buf[0])
	}
	disableRawMode()
}
