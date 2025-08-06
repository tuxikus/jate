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
//                                   Emacs                                   //
///////////////////////////////////////////////////////////////////////////////

func previousLine() {
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

func nextLine() {
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

func backwardChar() {
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

func forwardChar() {
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

func moveBeginningOfLine() {
	editor.cursorX = 0
}

func moveEndOfLine() {
	editor.cursorX = len(getCurrentRow().chars)
}

func backToIndentation() {
	currentRow := getCurrentRow()
	if currentRow != nil {
		editor.cursorX = getIndexOfFirstNonWhitespaceChar(currentRow)
	}
}

func forwardWord() {
	if indexOfWordEnd := getIndexOfWordEnd(); indexOfWordEnd != 0 {
		editor.cursorX = indexOfWordEnd
	} else {
		if indexOfNextWord := getIndexOfNextWord(); indexOfNextWord == 0 {
			editor.cursorX = 0
			editor.cursorY++
			forwardWord()
		} else {
			editor.cursorX = getIndexOfNextWord()
			forwardWord()
		}
	}
}

func backwardWord() {
	if indexOfWordBeginning := getIndexOfWordBeginning(); indexOfWordBeginning != editor.cursorX {
		editor.cursorX = indexOfWordBeginning
	} else {
		if editor.cursorX == 0 && editor.cursorY != 0 {
			editor.cursorY--
			editor.cursorX = len(editor.row[editor.cursorY].chars)
			backwardWord()
		} else {
			editor.cursorX = getIndexOfPreviousWord()
			backwardWord()
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
//                                     Vi                                    //
///////////////////////////////////////////////////////////////////////////////
