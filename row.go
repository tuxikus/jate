package main

func insertRow(at int, s string) {
	if at < 0 || at > editor.rows {
		return
	}

	// new row
	editor.row = append(editor.row, EditorRow{})

	// shift rows
	copy(editor.row[at+1:], editor.row[at:])

	editor.row[at].chars = []byte(s)
	editor.row[at].length = len(editor.row[at].chars)
	editor.row[at].render = nil
	editor.row[at].renderLength = 0

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
			row.render = append(row.render, '#')
			idx++
			for idx%TAB_WIDTH != 0 {
				row.render = append(row.render, '#')
				idx++
			}
		} else {
			row.render = append(row.render, char)
			idx++
		}
	}

	row.renderLength = len(row.render)
}
