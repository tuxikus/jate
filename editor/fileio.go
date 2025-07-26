package editor

import (
	"os"
)

func FileOpen(filename string) {
	if editor.filename == "" {
		editor.filename = filename
	} else {
		// if file already loaded, reset editor
		// TODO function
		editor.cursorX = 0
		editor.cursorY = 0
		editor.renderX = 0
		editor.rows = 0
		editor.rowOffset = 0
		editor.columnOffset = 0
		editor.row = nil
		editor.filename = ""
		editor.statusMessage = ""
		editor.fileModified = 0
		editor.syntax = nil

		editor.filename = filename
	}

	selectSyntax()

	content, err := os.ReadFile(filename)
	if err != nil {
		panicExit("open\n" + err.Error())
	}

	var contentAsBytes []byte

	for _, char := range content {
		if char == '\n' || char == 10 {
			insertRow(editor.rows, string(contentAsBytes))
			contentAsBytes = nil
			continue
		}
		contentAsBytes = append(contentAsBytes, char)
	}
	editor.fileModified = 0
}

func fileSave() {
	if editor.filename == "" {
		editor.filename = string(prompt("Save as: ", nil))
	}

	selectSyntax()

	file, err := os.OpenFile(editor.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panicExit("save " + err.Error())
	}
	defer file.Close()

	fileString := rowsToString()
	fileBytes := []byte(fileString)
	file.Write(fileBytes)
	setStatusMessage("%d bytes saved to disk", len(fileBytes))
	editor.fileModified = 0
}
