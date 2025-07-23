package main

func drawMessageBar(ab *AppendBuffer) {
	appendBufferAppend(ab, []byte("\x1b[K"))
	appendBufferAppend(ab, []byte(editor.statusMessage))
}
