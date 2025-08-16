// the append buffer (ab) is used to call the
// write syscall only once per 'frame'

package editor

type AppendBuffer struct {
	chars []byte
}

func appendBufferAppend(ab *AppendBuffer, chars []byte) {
	ab.chars = append(ab.chars, chars...)
}

func appendBufferAppendByte(ab *AppendBuffer, char byte) {
	ab.chars = append(ab.chars, char)
}
