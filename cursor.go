package main

func moveCursor(key int) {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}

	switch key {
	case KEY_ARROW_LEFT:
		if editor.cursorX != 0 {
			editor.cursorX--
		} else if editor.cursorY > 0 {
			editor.cursorY--
			editor.cursorX = editor.row[editor.cursorY].length
		}
	case KEY_ARROW_RIGHT:
		if row != nil && editor.cursorX < row.length {
			editor.cursorX++
		} else if row != nil && editor.cursorX == row.length {
			editor.cursorY++
			editor.cursorX = 0
		}
	case KEY_ARROW_UP:
		if editor.cursorY > 0 {
			editor.cursorY--
		}
	case KEY_ARROW_DOWN:
		if editor.cursorY < editor.rows {
			editor.cursorY++
		}
	}

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		// get the new row if y changed
		row = &editor.row[editor.cursorY]
	}

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
		}
	}
}
