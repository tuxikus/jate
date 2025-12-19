// this file maps keyboard inputs to
// editor actions like cursor movement

package editor

const (
	KEYMAP_NORMAL = iota
	KEYMAP_EMACS
	KEYMAP_VI
)

type keymap int
type keymapBindings map[int]func()

///////////////////////////////////////////////////////////////////////////////
//                                   Normal                                  //
///////////////////////////////////////////////////////////////////////////////

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
		KEY_M_COLON:     executeCommand,
		KEY_BACKSPACE:   deleteChar,
	}
}

///////////////////////////////////////////////////////////////////////////////
//                                     Vi                                    //
///////////////////////////////////////////////////////////////////////////////

var viKeymapBindings = keymapBindings{}

///////////////////////////////////////////////////////////////////////////////
//                                    init                                   //
///////////////////////////////////////////////////////////////////////////////

func init() {
	normalKeymapBindings = initNormalKeymapBindings()
}
