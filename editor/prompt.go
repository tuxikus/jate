package editor

import (
	"fmt"
	"strings"
	"unicode"
)

type PromptCallback func(input []byte, key int)

var lastCommand = ""
var lastCommandCandidates []string = nil
var lastCommandCandidateIdx = 0

func prompt(prompt string, promptCallback PromptCallback) []byte {
	buf := make([]byte, 0)

	for {
		setStatusMessage("%s%s", prompt, buf)
		draw()

		c := readKey()

		if c == KEY_BACKSPACE {
			if len(buf) != 0 {
				buf = buf[:len(buf)-1]
			}
		} else if c == '\x1b' {
			lastCommand = ""
			lastCommandCandidates = nil
			lastCommandCandidateIdx = 0

			setStatusMessage("")
			if promptCallback != nil {
				promptCallback(buf, c)
			}
			return nil
		} else if c == '\r' || c == '\n' {
			lastCommand = ""
			lastCommandCandidates = nil
			lastCommandCandidateIdx = 0

			if len(buf) != 0 {
				setStatusMessage("")
				return buf[:len(buf)]
			}
			if promptCallback != nil {
				promptCallback(buf, c)
			}
		} else if !unicode.IsControl(rune(c)) && c < 128 {
			buf = append(buf, byte(c))
			// completion
		} else if c == '\t' {
			completion(&buf)
		}

		if promptCallback != nil {
			promptCallback(buf, c)
		}
	}
}

func completion(buf *[]byte) {
	if lastCommand == "" {
		promptContent := string(*buf)

		// no command to complete
		if !strings.Contains(promptContent, ".") {
			return
		}

		// get command
		promptCommand := strings.Split(promptContent, ".")[0] // TODO check length

		for _, command := range commands.commands {
			if command.name == promptCommand {
				lastCommand = command.name
				lastCommandCandidates = command.candidates
				cycleCompletion(buf)
				break
			}
		}
	} else {
		cycleCompletion(buf)
	}

}

func cycleCompletion(buf *[]byte) {
	if len(lastCommandCandidates) == 0 {
		return
	}

	*buf = []byte(fmt.Sprintf("%s.%s", lastCommand, lastCommandCandidates[lastCommandCandidateIdx]))

	if lastCommandCandidateIdx == len(lastCommandCandidates)-1 {
		lastCommandCandidateIdx = 0
	} else {
		lastCommandCandidateIdx++
	}

}
