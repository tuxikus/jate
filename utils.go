package main

import (
	"fmt"
	"os"
)

func normalExit() {
	os.Stdin.Write([]byte("\x1b[2J")) // clear
	os.Stdin.Write([]byte("\x1b[H"))  // move cursor to 1 1
	disableRawMode()

	// for _, row := range editor.row {
	//	//fmt.Println(row.render)
	//	fmt.Println(row.chars)
	// }

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
