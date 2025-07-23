package main

import "strings"

var lastMatch = -1
var searchDirection = 1
var savedHlLine = 0
var savedHl []byte = nil

func searchCallback(query []byte, key int) {
	if savedHl != nil {
		setStatusMessage("Restoring hl")
		editor.row[savedHlLine].highlight = savedHl
		savedHl = nil
	}

	if key == '\r' || key == '\x1b' {
		lastMatch = -1
		searchDirection = 1
		return
	} else if key == KEY_ARROW_RIGHT || key == KEY_ARROW_DOWN {
		searchDirection = 1
	} else if key == KEY_ARROW_RIGHT || key == KEY_ARROW_UP {
		searchDirection = -1
	} else {
		lastMatch = -1
		searchDirection = 1
	}

	if lastMatch == -1 {
		searchDirection = 1
	}

	current := lastMatch

	for range editor.row {
		current += searchDirection

		if current == -1 {
			current = editor.rows - 1
		} else if current == editor.rows {
			current = 0
		}

		row := &editor.row[current]

		s := string(row.chars)
		match := strings.Index(s, string(query))
		if match != -1 {
			lastMatch = current
			editor.cursorY = current
			editor.cursorX = renderXtoCursorX(row, match)
			editor.rowOffset = editor.rows

			savedHlLine = current
			savedHl = row.highlight

			for i := match; i < match+len(query); i++ {
				row.highlight[i] = HL_SEARCH_MATCH
			}
			break
		}
	}
}

func search() {
	oldCursorX := editor.cursorX
	oldCursorY := editor.cursorY
	oldColumnOffset := editor.columnOffset
	oldRowOffset := editor.rowOffset

	if prompt("Search: ", searchCallback) == nil {
		editor.cursorX = oldCursorX
		editor.cursorY = oldCursorY
		editor.columnOffset = oldColumnOffset
		editor.rowOffset = oldRowOffset
	}
}
