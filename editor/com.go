package editor

import (
	"strings"
)

type Command struct {
	name       string
	candidates []string
}

type Commands struct {
	commands []Command
}

var commands = Commands{
	commands: []Command{
		{
			name:       "get",
			candidates: []string{"cursorX", "cursorY"},
		},
		{
			name:       "set",
			candidates: []string{"cursorX", "cursorY", "cursorZ"}, // cursorZ is used for testing
		},
		{
			name:       "open",
			candidates: getFiles("."),
		},
	},
}

func executeCommand() {
	command := string(prompt(":", nil))

	// TODO use reflection via reflect module

	// get
	// get editor variables
	// usage get.cursorX
	if strings.HasPrefix(command, "get") {
		editorVariable := ""
		if len(strings.Split(command, ".")) > 1 {
			editorVariable = strings.Split(command, ".")[1]
		} else {
			return
		}

		switch editorVariable {
		case "cursorX":
			setStatusMessage("%d", editor.cursorX)
		case "cursorY":
			setStatusMessage("%d", editor.cursorY)
		}
		// set
	} else if strings.HasPrefix(command, "set") {
		// TODO implement
		setStatusMessage("%s", "not implemented yet")
		// open
	} else if strings.HasPrefix(command, "open") {
		file := ""
		ext := ""
		if len(strings.Split(command, ".")) > 1 {
			file = strings.Split(command, ".")[1]
			ext = strings.Split(command, ".")[2]
		} else {
			return
		}
		FileOpen(file + "." + ext)
	}
}
