package main

func moveCursorUp() {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}

	if editor.cursorY > 0 {
		editor.cursorY--
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

func moveCursorDown() {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}
	if editor.cursorY < editor.rows {
		editor.cursorY++
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

func moveCursorLeft() {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}

	if editor.cursorX != 0 {
		editor.cursorX--
	} else if editor.cursorY > 0 {
		editor.cursorY--
		editor.cursorX = editor.row[editor.cursorY].length
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

func moveCursorRight() {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}

	if row != nil && editor.cursorX < row.length {
		editor.cursorX++
	} else if row != nil && editor.cursorX == row.length {
		editor.cursorY++
		editor.cursorX = 0
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
