// dirL is the standard file selector

package editor

import (
	"os"
)

var preDirLRows []EditorRow

func initDirLKeymapBindings() keymapBindings {
	return keymapBindings{
		13:      selectFile,
		KEY_C_P: previousLine,
		KEY_C_N: nextLine,
		'q':     closeDirL,
	}
}

func openDirL() {
	editor.keymapBindings = initDirLKeymapBindings()
	preDirLRows = window(append([]string{"=== dirL, the standard file selector ==="}, getCWDFiles()...))
}

func getCWDFiles() []string {
	filesStringList := make([]string, 0)

	dir, err := os.Getwd()
	if err != nil {
		return nil
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, file := range files {
		filesStringList = append(filesStringList, file.Name())
	}

	return filesStringList
}

func closeDirL() {
	editor.cursorX = 0
	editor.cursorY = 0

	editor.keymapBindings = emacsKeymapBindings
	editor.row = preDirLRows
}

func selectFile() {
	editor.keymapBindings = emacsKeymapBindings
	lineContent := string(getCurrentRow().chars)
	FileOpen(lineContent)
}
