package editor

import (
	"strings"
)

type Command struct {
	name       string
	aliases    []string
	candidates []string
}

type Commands struct {
	commands []Command
}

var commands = Commands{
	commands: []Command{
		{
			name:    "quit",
			aliases: []string{"exit", "q"},
		},
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
		{
			name:       "keyMode",
			candidates: []string{"emacs", "vi"},
		},
		{
			name: "help",
		},
	},
}

func executeCommand() {
	command := string(prompt(":", nil))

	// TODO use reflection via reflect module
	if strings.HasPrefix(command, "quit") {
		normalExit()
		// get
		// get editor variables
		// usage get.cursorX
	} else if strings.HasPrefix(command, "get") {
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
	} else if strings.HasPrefix(command, "keyMode") {
		mode := ""
		if len(strings.Split(command, ".")) > 1 {
			mode = strings.Split(command, ".")[1]

			switch mode {
			case "emacs":
				editor.keyBindingMode = KEY_BINDING_MODE_EMACS
			case "vi":
				editor.keyBindingMode = KEY_BINDING_MODE_VI
			}
		} else {
			return
		}
	} else if strings.HasPrefix(command, "help") {
		setStatusMessage("%s", "commands: quit, get, set, open, keyMode, help")
	}
}
