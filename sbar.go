package main

import "fmt"

func drawStatusBar(ab *AppendBuffer) {
	appendBufferAppend(ab, []byte("\x1b[7m"))

	left := fmt.Sprintf("File: %s", editor.filename)
	if editor.fileModified != 0 {
		left += " -modified-"
	}

	right := fmt.Sprintf("Lines: %d", editor.rows)

	appendBufferAppend(ab, []byte(left))
	for range editor.screenColumns - len(left) - len(right) {
		appendBufferAppend(ab, []byte(" "))
	}
	appendBufferAppend(ab, []byte(right))
	appendBufferAppend(ab, []byte("\x1b[m"))
	appendBufferAppend(ab, []byte("\r\n"))
}
