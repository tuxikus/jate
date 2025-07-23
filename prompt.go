package main

import "unicode"

type PromptCallback func(input []byte, key int)

func prompt(prompt string, promptCallback PromptCallback) []byte {
	// bufSize := 128
	// buflen := 0
	buf := make([]byte, 0)

	for {
		setStatusMessage("%s%s", prompt, buf)
		refreshScreen()

		c := readKey()

		if c == KEY_BACKSPACE {
			if len(buf) != 0 {
				buf = buf[:len(buf)-1]
			}
		} else if c == '\x1b' {
			setStatusMessage("")
			if promptCallback != nil {
				promptCallback(buf, c)
			}
			return nil
		} else if c == '\r' || c == '\n' {
			if len(buf) != 0 {
				setStatusMessage("")
				return buf[:len(buf)]
			}
			if promptCallback != nil {
				promptCallback(buf, c)
			}
		} else if !unicode.IsControl(rune(c)) && c < 128 {
			buf = append(buf, byte(c))
		}

		if promptCallback != nil {
			promptCallback(buf, c)
		}
	}
}
