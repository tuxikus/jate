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

var normalKeymapBindings = keymapBindings{
	KEY_ARROW_UP:    moveCursorUp,
	KEY_ARROW_RIGHT: moveCursorRight,
	KEY_ARROW_DOWN:  moveCursorDown,
	KEY_ARROW_LEFT:  moveCursorLeft,
	KEY_HOME:        moveCursorToBeginning,
	KEY_END:         moveCursorToEnd,
	KEY_M_COLON:     executeCommand,
	KEY_BACKSPACE:   deleteChar,
}

///////////////////////////////////////////////////////////////////////////////
//                                   Emacs                                   //
///////////////////////////////////////////////////////////////////////////////

var emacsKeymapBindings = keymapBindings{
	KEY_C_X: func() {
		if fn := prefix(KEY_C_X); fn != nil {
			fn()
		}
	},
	KEY_C_C: func() {
		if fn := prefix(KEY_C_C); fn != nil {
			fn()
		}
	},
	KEY_C_P:       previousLine,
	KEY_C_F:       forwardChar,
	KEY_C_N:       nextLine,
	KEY_C_B:       backwardChar,
	KEY_C_A:       moveBeginningOfLine,
	KEY_C_E:       moveEndOfLine,
	KEY_C_K:       killLine,
	KEY_M_LOWER_F: forwardWord,
	KEY_M_LOWER_B: backwardWord,
	KEY_M_LOWER_M: backToIndentation,
	KEY_M_COLON:   executeCommand,
	KEY_BACKSPACE: deleteChar,
}

// prefix implementation
// return the action as a function
func prefix(prefixKey int) func() {
	c := readKey()

	switch prefixKey {
	case KEY_C_X:
		switch c {
		case KEY_C_C:
			return normalExit
		case KEY_C_S:
			return fileSave
		default:
			return nil
		}
	default:
		return nil
	}
}

///////////////////////////////////////////////////////////////////////////////
//                                     Vi                                    //
///////////////////////////////////////////////////////////////////////////////

var viKeymapBindings = keymapBindings{}
