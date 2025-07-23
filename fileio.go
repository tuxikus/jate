package main

import (
	"fmt"
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

	file, err := os.OpenFile(editor.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panicExit("save " + err.Error())
	}
	defer file.Close()

	fileBytes := []byte(rowsToString())
	file.Write(fileBytes)
	setStatusMessage(fmt.Sprintf("%d bytes saved to disk!", len(fileBytes)))
	editor.fileModified = 0
}
