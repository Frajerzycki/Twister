package test

import (
	"bytes"
)

// TestBuffer is used for tests.
type TestBuffer struct{ buffer *bytes.Buffer }

func (buffer TestBuffer) Write(data []byte) (int, error) {
	return buffer.buffer.Write(data)
}

func (buffer TestBuffer) Read(destiny []byte) (int, error) {
	return buffer.buffer.Read(destiny)
}

// Close closes TestBuffer.
func (buffer TestBuffer) Close() error {
	return nil
}

// Bytes returns data stored in TestBuffer.
func (buffer TestBuffer) Bytes() []byte {
	return buffer.buffer.Bytes()
}

// NewTestBuffer creates new TestBuffer.
func NewTestBuffer() TestBuffer {
	return TestBuffer{buffer: new(bytes.Buffer)}
}
