package main

import (
	"fmt"
	"os"
	"unicode"
)

func draw() {
	var appendBuffer AppendBuffer

	scroll()

	// hide the cursor
	appendBufferAppend(&appendBuffer, []byte("\x1b?25l"))

	// clear the screen
	appendBufferAppend(&appendBuffer, []byte("\x1b[2J"))

	// reposition the cursor to the beginning
	// H: VT100 cursor position
	// [10;10H move cursor to row 10 and column 10
	// default is 1;1
	appendBufferAppend(&appendBuffer, []byte("\x1b[H"))

	drawRows(&appendBuffer)       // screenrows - 2
	drawStatusBar(&appendBuffer)  // screenrows - 1
	drawMessageBar(&appendBuffer) // screenrows

	cursorVt100 := fmt.Sprintf("\x1b[%d;%dH", editor.cursorY-editor.rowOffset+1, editor.renderX-editor.columnOffset+1)
	appendBufferAppend(&appendBuffer, []byte(cursorVt100))

	// show the cursor
	appendBufferAppend(&appendBuffer, []byte("\x1b[?25h"))

	os.Stdin.Write(appendBuffer.chars) // the only write call per refresh
}

// visual part of the editor
func drawRows(ab *AppendBuffer) {
	for y := range editor.screenRows {
		filerow := y + editor.rowOffset
		// print ~ after the file content
		if filerow >= editor.rows {
			// only display the welcome message if no file is loaded
			if editor.rows == 0 && y == editor.screenRows/2 {
				// message
				// draw tilde at start of line
				appendBufferAppend(ab, []byte("~"))
				// draw line 1
				messageLine1 := "jate - just another text editor"
				padding := ((editor.screenColumns - len(messageLine1)) / 2) - 1 // -1 = tilde
				for padding > 0 {
					appendBufferAppend(ab, []byte(" "))
					padding--
				}
				appendBufferAppend(ab, []byte(messageLine1))

				// fill line 1 so that line 2 is centered
				padding = (editor.screenColumns - len(messageLine1)) / 2
				for padding > 0 {
					appendBufferAppend(ab, []byte(" "))
					padding--
				}

				// draw tilde at start of line
				appendBufferAppend(ab, []byte("~"))
				// draw line 2
				messageLine2 := fmt.Sprintf("Version: %s", VERSION)
				padding = ((editor.screenColumns - len(messageLine2)) / 2) - 1 // -1 = tilde
				for padding > 0 {
					appendBufferAppend(ab, []byte(" "))
					padding--
				}
				appendBufferAppend(ab, []byte(messageLine2))
			} else {
				appendBufferAppend(ab, []byte("~"))
			}
		} else {
			drawStart := min(editor.columnOffset, len(editor.row[filerow].render))
			drawEnd := min(drawStart+editor.screenColumns, len(editor.row[filerow].render))

			// no color
			//appendBufferAppend(ab, editor.row[filerow].render[drawStart:drawEnd])

			// some color
			rowChars := editor.row[filerow].render[drawStart:drawEnd]
			hl := editor.row[filerow].highlight[drawStart:drawEnd]
			currentColor := -1
			for i := range drawEnd {
				if unicode.IsControl(rune(rowChars[i])) {
					sym := make([]byte, 0)
					if rowChars[i] <= 26 {
						sym = append(sym, '@', rowChars[i])
					} else {
						sym = append(sym, '?')
					}
					appendBufferAppend(ab, []byte("\x1b[7m"))
					appendBufferAppend(ab, sym)
					appendBufferAppend(ab, []byte("\x1b[m"))
					if currentColor != -1 {
						colorString := fmt.Sprintf("\x1b[%dm", currentColor)
						appendBufferAppend(ab, []byte(colorString))
					}
				} else if hl[i] == HL_NORMAL {
					if currentColor != -1 {
						appendBufferAppend(ab, []byte("\x1b[39m"))
						currentColor = -1
					}
					appendBufferAppendByte(ab, rowChars[i])
				} else {
					color := syntaxToColor(int(hl[i]))
					if color != currentColor {
						currentColor = color
						colorString := fmt.Sprintf("\x1b[%dm", color)
						appendBufferAppend(ab, []byte(colorString))
					}
					appendBufferAppendByte(ab, rowChars[i])
				}
			}
			appendBufferAppend(ab, []byte("\x1b[39m"))
		}

		appendBufferAppend(ab, []byte("\x1b[K"))
		appendBufferAppend(ab, []byte("\r\n"))
	}
}
