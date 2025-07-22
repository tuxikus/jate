package main

import (
	"fmt"
	"os"
)

func FileOpen(filename string) {
	editor.filename = filename

	content, err := os.ReadFile(filename)
	if err != nil {
		PanicExit("open")
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

func FileSave() {
	if editor.filename == "" {
		editor.filename = string(Prompt("Save as: ", nil))
	}

	// TODO overwrite the open file
	//file, err := os.Open("./temp.txt")
	file, err := os.Create("./temp.txt")
	if err != nil {
		PanicExit("save " + err.Error())
	}
	defer file.Close()

	fileBytes := []byte(rowsToString())
	file.Write(fileBytes)
	setStatusMessage(fmt.Sprintf("%d bytes saved to disk!", len(fileBytes)))
	editor.fileModified = 0
}
