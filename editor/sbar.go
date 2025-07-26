package editor

import (
	"fmt"
	"time"
)

func drawStatusBar(ab *AppendBuffer) {
	appendBufferAppend(ab, []byte("\x1b[7m"))

	var fType []byte
	var fName []byte

	if editor.syntax != nil && editor.syntax.fileType != nil {
		fType = []byte("[" + string(editor.syntax.fileType) + "]")
	} else {
		fType = []byte("")
	}

	if editor.filename != "" {
		fName = []byte(editor.filename)
	} else {
		fName = []byte("-")
	}

	left := fmt.Sprintf(" %s File: %s Lines: %d:%d", fType, fName, editor.rows, editor.cursorY+1)
	if editor.fileModified != 0 {
		left += " -modified-"
	}

	t := time.Now().Format("15:04")
	right := fmt.Sprintf("%s  ", t)

	appendBufferAppend(ab, []byte(left))
	for range editor.screenColumns - len(left) - len(right) {
		appendBufferAppend(ab, []byte(" "))
	}
	appendBufferAppend(ab, []byte(right))
	appendBufferAppend(ab, []byte("\x1b[m"))
	appendBufferAppend(ab, []byte("\r\n"))
}
