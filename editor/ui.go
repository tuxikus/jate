// file for ui elements

package editor

func blankWindow() {
	editor.cursorX = 0

	for range editor.screenRows {
		insertRow(0, "")
	}
}

func window(content []string) []EditorRow {
	old := editor.row

	blankWindow()

	for i := range content {
		editor.row[i].chars = []byte(content[i])
		updateRow(&editor.row[i])
	}

	return old
}
