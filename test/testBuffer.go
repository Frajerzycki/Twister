package test

import (
	"bytes"
)

type TestBuffer struct{ buffer *bytes.Buffer }

func (buffer TestBuffer) Write(data []byte) (int, error) {
	return buffer.buffer.Write(data)
}

func (buffer TestBuffer) Read(destiny []byte) (int, error) {
	return buffer.buffer.Read(destiny)
}

func (buffer TestBuffer) Close() error {
	return nil
}

func (buffer TestBuffer) Bytes() []byte {
	return buffer.buffer.Bytes()
}

func NewTestBuffer() TestBuffer {
	return TestBuffer{buffer: new(bytes.Buffer)}
}
