package editor

///////////////////////////////////////////////////////////////////////////////
//                                   Normal                                  //
///////////////////////////////////////////////////////////////////////////////

func moveCursorUp() {
	if editor.cursorY > 0 {
		editor.cursorY--
	}

	currentRow := getCurrentRow()

	// check if cursor is past the row length
	if currentRow != nil {
		if editor.cursorX > currentRow.length {
			editor.cursorX = currentRow.length
		}
	}

}

func moveCursorDown() {
	if editor.cursorY < editor.rows {
		editor.cursorY++
	}

	currentRow := getCurrentRow()

	// check if cursor is past the row length
	if currentRow != nil {
		if editor.cursorX > currentRow.length {
			editor.cursorX = currentRow.length
		}
	}
}

func moveCursorLeft() {
	if editor.cursorX != 0 {
		editor.cursorX--
	} else if editor.cursorY > 0 {
		editor.cursorY--
		editor.cursorX = len(editor.row[editor.cursorY].chars)
	}

	currentRow := getCurrentRow()

	// check if cursor is past the row length
	if currentRow != nil {
		if editor.cursorX > currentRow.length {
			editor.cursorX = currentRow.length
		}
	}
}

func moveCursorRight() {
	currentRow := getCurrentRow()

	if currentRow != nil && editor.cursorX < currentRow.length {
		editor.cursorX++
	} else if currentRow != nil && editor.cursorX == currentRow.length {
		editor.cursorY++
		editor.cursorX = 0
	}

	currentRow = getCurrentRow()

	// check if cursor is past the row length
	if currentRow != nil {
		if editor.cursorX > currentRow.length {
			editor.cursorX = currentRow.length
		}
	}
}

func moveCursorToBeginning() {
	editor.cursorX = 0
}

func moveCursorToEnd() {
	editor.cursorX = len(getCurrentRow().chars)
}

///////////////////////////////////////////////////////////////////////////////
//                                     Vi                                    //
///////////////////////////////////////////////////////////////////////////////
