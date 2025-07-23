package main

func deleteRow(at int) {
	if at < 0 || at >= editor.rows {
		return
	}

	// copy(dst, src)
	copy(editor.row[at:], editor.row[at+1:])

	for i := at; i < editor.rows-1; i++ {
		editor.row[i].idx--
	}

	editor.rows--
	editor.fileModified++
}

func rowAppendString(row *EditorRow, s string) {
	//row.chars = make([]byte, len(row.chars) + len(s))
	//copy(row.chars[:len(row.chars)], s)

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
	for _, row := range editor.row {
		s += string(row.chars)
		s += "\n"
	}
	return s
}

func insertRow(at int, s string) {
	if at < 0 || at > editor.rows {
		return
	}

	// new row
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

func updateRow(row *EditorRow) {
	// count tabs
	tabs := 0
	for _, char := range row.chars {
		if char == '\t' {
			tabs++
		}
	}

	size := len(row.chars) + tabs*(TAB_WIDTH-1) + 1
	row.render = make([]byte, 0, size)

	idx := 0
	for _, char := range row.chars {
		if char == '\t' {
			row.render = append(row.render, ' ')
			idx++
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
