package editor

// used to call 'write' only once per refresh
type AppendBuffer struct {
	chars []byte
}

func appendBufferAppend(ab *AppendBuffer, chars []byte) {
	ab.chars = append(ab.chars, chars...)
}

func appendBufferAppendByte(ab *AppendBuffer, char byte) {
	ab.chars = append(ab.chars, char)
}
