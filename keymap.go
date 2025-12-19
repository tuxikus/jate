// this file maps keyboard inputs to
// editor actions like cursor movement

package main

const (
	KEYMAP_NORMAL = iota
)

type keymap int
type keymapBindings map[int]func()

var normalKeymapBindings keymapBindings

func initNormalKeymapBindings() keymapBindings {
	return keymapBindings{
		13:              insertNewLine, // 13 => Carriage Return (Enter)
		KEY_ARROW_UP:    moveCursorUp,
		KEY_ARROW_RIGHT: moveCursorRight,
		KEY_ARROW_DOWN:  moveCursorDown,
		KEY_ARROW_LEFT:  moveCursorLeft,
		KEY_HOME:        moveCursorToBeginning,
		KEY_END:         moveCursorToEnd,
		KEY_BACKSPACE:   deleteChar,
		KEY_C_C:         normalExit,
		KEY_C_S:         fileSave,
		KEY_C_F:         search,
	}
}

func init() {
	normalKeymapBindings = initNormalKeymapBindings()
}
