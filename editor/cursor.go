package editor

///////////////////////////////////////////////////////////////////////////////
//                                   Normal                                  //
///////////////////////////////////////////////////////////////////////////////

func moveCursorUp() {
	if editor.cursorY > 0 {
		editor.cursorY--
	}

	row := getCurrentRow()

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
		}
	}

}

func moveCursorDown() {
	if editor.cursorY < editor.rows {
		editor.cursorY++
	}

	row := getCurrentRow()

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
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

	row := getCurrentRow()

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
		}
	}
}

func moveCursorRight() {
	row := getCurrentRow()

	if row != nil && editor.cursorX < row.length {
		editor.cursorX++
	} else if row != nil && editor.cursorX == row.length {
		editor.cursorY++
		editor.cursorX = 0
	}

	row = getCurrentRow()

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
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
//                                   Emacs                                   //
///////////////////////////////////////////////////////////////////////////////

func previousLine() {
	if editor.cursorY > 0 {
		editor.cursorY--
	}

	row := getCurrentRow()

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
		}
	}

}

func nextLine() {
	if editor.cursorY < editor.rows {
		editor.cursorY++
	}

	row := getCurrentRow()

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
		}
	}
}

func forwardChar() {
	if editor.cursorX != 0 {
		editor.cursorX--
	} else if editor.cursorY > 0 {
		editor.cursorY--
		editor.cursorX = len(editor.row[editor.cursorY].chars)
	}

	row := getCurrentRow()

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
		}
	}
}

func backwardChar() {
	row := getCurrentRow()

	if row != nil && editor.cursorX < row.length {
		editor.cursorX++
	} else if row != nil && editor.cursorX == row.length {
		editor.cursorY++
		editor.cursorX = 0
	}

	row = getCurrentRow()

	// check if cursor is past the row length
	if row != nil {
		if editor.cursorX > row.length {
			editor.cursorX = row.length
		}
	}
}

func moveBeginningOfLine() {
	editor.cursorX = 0
}

func moveEndOfLine() {
	editor.cursorX = len(getCurrentRow().chars)
}

func backToIndentation() {
	row := getCurrentRow()
	if row != nil {
		editor.cursorX = getIndexOfFirstNonWhitespaceChar(row)
	}
}

///////////////////////////////////////////////////////////////////////////////
//                                     Vi                                    //
///////////////////////////////////////////////////////////////////////////////

func moveCursorLeftVi() {
	if editor.cursorX != 0 {
		editor.cursorX--
	}
}

func moveCursorRightVi() {
	if editor.cursorX < len(editor.row[editor.cursorY].chars) {
		editor.cursorX++
	}
}

func moveCursorDownVi() {
	// on last line
	if editor.cursorY+1 >= editor.rows {
		return
	}

	if editor.cursorY < editor.rows {
		if editor.cursorX > len(editor.row[editor.cursorY+1].chars) {
			editor.cursorX = len(editor.row[editor.cursorY+1].chars) - 1
		}
		editor.cursorY++
	}
}

func moveCursorUpVi() {
	if editor.cursorY > 0 {
		if editor.cursorX > len(editor.row[editor.cursorY-1].chars) {
			editor.cursorX = len(editor.row[editor.cursorY-1].chars) - 1
		}
		editor.cursorY--
	}
}

func moveCursorToIndentation() {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}

	x := 0
	for _, char := range row.chars {
		if char == ' ' || char == '\t' {
			x++
		} else {
			break
		}
	}

	editor.cursorX = x
}

func moveCursorWordForward() {
	var row *EditorRow
	if editor.cursorY >= editor.rows {
		return
	} else {
		row = &editor.row[editor.cursorY]
	}

	inWord := false

	// if cursor at the end of line move to the next line with chars
	if editor.cursorX >= len(row.chars) {
		editor.cursorX = 0
		editor.cursorY++

		if editor.cursorY >= editor.rows {
			return
		}

		row = &editor.row[editor.cursorY]

		// find next row with non symbol chars
		for {
			if len(row.chars) != 0 && rowContainsLetterOrDigit(row) {
				return
			} else {
				editor.cursorY++
				row = &editor.row[editor.cursorY]
			}
		}

	}

	if !isSymbol(row.chars[editor.cursorX]) {
		inWord = true
	}

	for i := editor.cursorX; i < len(editor.row[editor.cursorY].chars); i++ {
		// if in word move the cursor to the end of the word
		if inWord {
			if isSymbol(row.chars[editor.cursorX]) {
				return
			}
			editor.cursorX++
			// if not in word move the cursor to the next word and
			// set inWord to true and move to the end of this word
		} else {
			editor.cursorX++

			if editor.cursorX >= len(row.chars) {
				editor.cursorX = 0
				editor.cursorY++
				row = &editor.row[editor.cursorY]
				for len(row.chars) == 0 {
					editor.cursorY++
					row = &editor.row[editor.cursorY]
				}
			}

			if !isSymbol(row.chars[editor.cursorX]) {
				inWord = true
			}
		}
	}
}

func moveCursorWordBackward() {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}

	// if cursor at the beginning of line move to the previous line with chars
	if editor.cursorX == 0 {
		editor.cursorY--

		if editor.cursorY < 0 {
			return
		}

		row = &editor.row[editor.cursorY]
		editor.cursorX = len(row.chars)

		for len(row.chars) == 0 {
			editor.cursorY--
			row = &editor.row[editor.cursorY]
			editor.cursorX = len(row.chars)
		}
	}

	inWord := false
	toNextWord := false

	if editor.cursorX >= len(row.chars) {
		editor.cursorX = len(row.chars) - 1
	}

	if !isSymbol(row.chars[editor.cursorX]) {
		inWord = true

		if editor.cursorX-1 <= 0 {
			return
		}

		if isSymbol(row.chars[editor.cursorX-1]) {
			setStatusMessage("to next word")
			toNextWord = true
			inWord = false
		}
	}

	for {
		if editor.cursorX-1 <= 0 {
			editor.cursorX = 0
			return
		}

		if inWord {
			if isSymbol(row.chars[editor.cursorX-1]) {
				return
			}
			editor.cursorX--
		} else if toNextWord {
			editor.cursorX--
			if !isSymbol(row.chars[editor.cursorX]) {
				inWord = true
				toNextWord = false
			}
		} else {
			editor.cursorX--
			if !isSymbol(row.chars[editor.cursorX]) {
				inWord = true
			}
		}
	}
}
