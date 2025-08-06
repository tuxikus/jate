package editor

import (
	"fmt"
	"os"
	"slices"
)

func normalExit() {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H"))  // move cursor to 1 1
	disableRawMode()

	for _, row := range editor.row {
		//fmt.Println(row.render)
		fmt.Println(row.chars)
	}

	os.Exit(0)
}

func panicExit(message string) {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H"))  // move cursor to 1 1
	disableRawMode()
	fmt.Println(message)
	os.Exit(1)
}

func renderXtoCursorX(row *EditorRow, renderX int) int {
	currentRenderX := 0

	for cursorX := 0; cursorX < len(row.chars); cursorX++ {
		if row.chars[cursorX] == '\t' {
			currentRenderX += (TAB_WIDTH - 1) - (currentRenderX % TAB_WIDTH)
		}
		currentRenderX++

		if currentRenderX > renderX {
			return cursorX
		}
	}

	return len(row.chars)
}

func cursorXToRenderX(row *EditorRow, cursorX int) int {
	renderX := 0
	for i := range cursorX {
		if row.chars[i] == '\t' {
			// how many columns right to the last tab
			renderX += TAB_WIDTH - 1 - (renderX % TAB_WIDTH)
		}
		renderX++
	}

	return renderX
}

func isSymbol(c byte) bool {
	symbols := []byte{
		'!',
		'"',
		'#',
		'$',
		'%',
		'%',
		'\'',
		'(',
		')',
		'*',
		'+',
		',',
		'-',
		'.',
		'/',
		':',
		';',
		'<',
		'=',
		'>',
		'?',
		'@',
		'[',
		'\\',
		']',
		'^',
		'_',
		'`',
		'{',
		'|',
		'}',
		'~',
		' ',
		'\t',
		'\r',
	}
	if slices.Contains(symbols, c) {
		return true
	} else {
		return false
	}

}

func rowContainsLetterOrDigit(row *EditorRow) bool {
	if len(row.chars) == 0 {
		return false
	}

	for _, char := range row.chars {
		if !isSymbol(char) {
			return true
		}
	}

	return false
}

func getFiles(directory string) []string {
	var files []string = nil

	dir, err := os.Open(directory)
	if err != nil {
		panicExit("getFiles\n" + err.Error())
	}

	f, err := dir.ReadDir(-1)
	if err != nil {
		panicExit("getFiles\n" + err.Error())
	}

	for _, file := range f {
		files = append(files, file.Name())
	}

	return files
}
