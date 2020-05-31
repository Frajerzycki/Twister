package files

import (
	"errors"
	"io"
)

// ChunkedReader reads from Reader in chunks, but use as much data as needed.
type ChunkedReader struct {
	chunkSize         int
	chunk             []byte
	chunkMemory       []byte
	isReaderAvailable bool
	Reader            io.Reader
}

// NewChunkedReader creates new ChunkedReader.
func NewChunkedReader(reader io.Reader, chunkSize int) (*ChunkedReader, error) {
	if chunkSize < 1 {
		return nil, errors.New("Chunk size has to be positive.")
	}
	return &ChunkedReader{Reader: reader, chunkSize: chunkSize, chunkMemory: make([]byte, chunkSize), chunk: nil, isReaderAvailable: true}, nil
}

func removeFirstBytes(bytes []byte, size int) []byte {
	if len(bytes) > size {
		return bytes[size:]
	}
	return nil
}

func (chunkedReader *ChunkedReader) Read(destiny []byte) (int, error) {
	destinyFirstLength := len(destiny)
	for len(destiny) > 0 {
		if len(chunkedReader.chunk) == 0 && chunkedReader.isReaderAvailable {
			bytesRead, err := io.ReadFull(chunkedReader.Reader, chunkedReader.chunkMemory)
			chunkedReader.chunk = chunkedReader.chunkMemory[:bytesRead]
			if err != nil {
				chunkedReader.isReaderAvailable = false
			}
		}

		bytesCopied := copy(destiny, chunkedReader.chunk)
		destiny = removeFirstBytes(destiny, bytesCopied)
		chunkedReader.chunk = removeFirstBytes(chunkedReader.chunk, bytesCopied)
		if len(destiny) > 0 && !chunkedReader.isReaderAvailable {
			return destinyFirstLength - len(destiny), io.EOF
		}
	}

	return destinyFirstLength, nil
}
