package editor

import "fmt"

func drawStatusBar(ab *AppendBuffer) {
	appendBufferAppend(ab, []byte("\x1b[7m"))

	var fType []byte

	if editor.syntax != nil && editor.syntax.fileType != nil {
		fType = editor.syntax.fileType
	} else {
		fType = []byte("-")
	}

	left := fmt.Sprintf("Type: %s File: %s", fType, editor.filename)
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
