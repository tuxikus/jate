package main

import (
	"os"
)

func fileOpen(filename string) {
	editor.filename = filename

	selectSyntax()

	content, err := os.ReadFile(filename)
	if err != nil {
		panicExit("open")
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
