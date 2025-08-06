package editor

// TODO: write test
func getIndexOfPreviousWord() int {
	currentRow := getCurrentRow()

	for i := editor.cursorX - 1; i > 0; i-- {
		if !isSymbol(currentRow.chars[i]) {
			return i
		}
	}

	return -1
}

// TODO: write test
func getIndexOfWordBeginning() int {
	currentRow := getCurrentRow()

	for i := editor.cursorX - 1; i >= 0; i-- {
		if currentChar := currentRow.chars[i]; isSymbol(currentChar) {
			return i + 1
		}
	}

	return 0
}

// TODO: write test
// only works if cursor is on whitespace or symbol
func getIndexOfNextWord() int {
	currentRow := getCurrentRow()

	for i := editor.cursorX; i < len(currentRow.chars); i++ {
		if !isSymbol(currentRow.chars[i]) {
			return i
		}
	}

	return 0
}

// TODO: write test
// if cursor is in a word return the index of the end of this word
// if cursor is not in a word return -1
func getIndexOfWordEnd() int {
	currentRow := getCurrentRow()

	// no row: end of file or empty file
	if currentRow == nil {
		return -1
	}

	if editor.cursorX >= len(currentRow.chars) {
		return 0
	}

	if isSymbol(currentRow.chars[editor.cursorX]) {
		return 0
	}

	for i := editor.cursorX; i < len(currentRow.chars); i++ {
		// TODO: fix if symbols after last word
		// last word of line
		if currentRowLength := len(currentRow.chars); i == currentRowLength-1 {
			return currentRowLength
		}

		if currentChar := currentRow.chars[i]; isSymbol(currentChar) {
			return i
		}
	}

	return 0
}

// TODO: write test
func getCurrentChar() byte {
	return getCurrentRow().chars[editor.cursorX]
}

// TODO: write test
func getCurrentRow() *EditorRow {
	var row *EditorRow

	if editor.cursorY >= editor.rows {
		row = nil
	} else {
		row = &editor.row[editor.cursorY]
	}

	return row
}

// TODO: test
func getIndexOfFirstNonWhitespaceChar(row *EditorRow) int {
	for i, char := range row.chars {
		if char != ' ' && char != '\t' {
			return i
		}
	}

	return len(row.chars)
}

// deletes a row from editor.row
func rowDelete(at int) {
	if at < 0 || at >= editor.rows {
		return
	}

	// copy(dst, src)
	// copy all rows below at to the index of at
	copy(editor.row[at:], editor.row[at+1:])

	for i := at; i < editor.rows-1; i++ {
		editor.row[i].idx--
	}

	editor.rows--
	editor.fileModified++
}

// delete the content of a row,
// if the line has no content delete the row
func rowDeleteContent(at int) {
	if at < 0 || at >= editor.rows {
		return
	}

	row := &editor.row[at]

	if len(row.chars) == 0 {
		rowDelete(at)
		return
	}

	row.chars = make([]byte, 0)
	updateRow(row)
}

func rowAppendString(row *EditorRow, s string) {
	row.chars = append(row.chars, s...)
	row.length = len(row.chars)
	updateRow(row)
	editor.fileModified++
}

func rowDeleteChar(row *EditorRow, at int) {
	if at < 0 || at > row.length {
		return
	}

	// copy(dst, src)
	copy(row.chars[at:], row.chars[at+1:])
	row.chars = row.chars[:len(row.chars)-1]
	row.length--
	updateRow(row)
	editor.fileModified++
}

func rowInsertChar(row *EditorRow, at int, char byte) {
	if at < 0 || at > row.length {
		at = row.length
	}

	row.chars = append(row.chars, 0)       // add one char to make room for new char
	copy(row.chars[at+1:], row.chars[at:]) // shift all chars from at to the right
	row.chars[at] = char                   // add the char

	row.length++
	updateRow(row)
}

func rowsToString() string {
	s := ""
	for i := range editor.rows {
		s += string(editor.row[i].chars) + "\n"
	}

	return s
}

func insertRow(at int, s string) {
	if at < 0 || at > editor.rows {
		return
	}

	// new empty row
	editor.row = append(editor.row, EditorRow{})

	// shift rows
	copy(editor.row[at+1:], editor.row[at:])

	for i := at + 1; i <= editor.rows; i++ {
		editor.row[i].idx++
	}

	editor.row[at].idx = at

	editor.row[at].chars = []byte(s)
	editor.row[at].length = len(editor.row[at].chars)
	editor.row[at].render = nil
	editor.row[at].highlight = nil
	editor.row[at].renderLength = 0
	editor.row[at].hlOpenComment = 0

	updateRow(&editor.row[at])

	editor.rows++
	editor.fileModified++
}

// build the row.render from row.chars
func updateRow(row *EditorRow) {
	// count tabs
	tabs := 0
	for _, char := range row.chars {
		if char == '\t' {
			tabs++
		}
	}

	// TAB_WIDTH - 1 -> /t already a char
	size := len(row.chars) + tabs*(TAB_WIDTH-1)
	row.render = make([]byte, 0, size)

	idx := 0
	for _, char := range row.chars {
		if char == '\t' {
			row.render = append(row.render, ' ')
			idx++
			// if char is a tab check idx and add needed spaces to fill the tab
			for idx%TAB_WIDTH != 0 {
				row.render = append(row.render, ' ')
				idx++
			}
		} else {
			row.render = append(row.render, char)
			idx++
		}
	}

	row.renderLength = len(row.render)

	updateSyntax(row)
}
