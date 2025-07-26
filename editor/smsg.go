package editor

import (
	"fmt"
)

func setStatusMessage(format string, a ...interface{}) {
	editor.setStatusMessage(fmt.Sprintf(format, a...))
}
