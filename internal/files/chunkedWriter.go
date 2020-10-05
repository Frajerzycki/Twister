package files

import (
	"errors"
	"io"
)

// ChunkedWriter writes to writer when chunk is filled.
type ChunkedWriter struct {
	chunkSize   int
	chunkMemory []byte
	chunk       []byte
	writer      io.Writer
}

// NewChunkedWriter creates new ChunkedWriter.
func NewChunkedWriter(writer io.Writer, chunkSize int) (*ChunkedWriter, error) {
	if chunkSize < 1 {
		return nil, errors.New("Chunk size has to be positive.")
	}

	chunkMemory := make([]byte, chunkSize)
	return &ChunkedWriter{chunkSize: chunkSize, chunkMemory: chunkMemory, chunk: chunkMemory, writer: writer}, nil
}

func copyAndRemoveBytes(to, from []byte) ([]byte, []byte) {
	bytesCopied := copy(to, from)
	return removeFirstBytes(to, bytesCopied), removeFirstBytes(from, bytesCopied)
}

func (writer *ChunkedWriter) writeAndClearChunk() error {
	writer.chunk = writer.chunkMemory
	_, err := writer.writer.Write(writer.chunk)
	if err != nil {
		return err
	}
	return nil
}

func (writer *ChunkedWriter) Write(data []byte) (int, error) {
	dataFirstLength := len(data)

	writer.chunk, data = copyAndRemoveBytes(writer.chunk, data)
	for len(data) > 0 {
		err := writer.writeAndClearChunk()
		if err != nil {
			return dataFirstLength - len(data), err
		}
		writer.chunk, data = copyAndRemoveBytes(writer.chunk, data)
	}
	return dataFirstLength, nil
}

// Close closes ChunkedWriter and writes unwritten data.
func (writer *ChunkedWriter) Close() error {
	writer.chunkMemory = writer.chunkMemory[:writer.chunkSize-len(writer.chunk)]
	return writer.writeAndClearChunk()
}
