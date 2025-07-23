package main

import "fmt"

func setStatusMessage(format string, a ...interface{}) {
	editor.statusMessage = fmt.Sprintf(format, a...)
}
