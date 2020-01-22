package functionality

import "io"

type DataBlocks struct {
	Rest   int
	Blocks [][]byte
}

func ByteArrayToDataBlocks(bytes []byte, blockLength int, reader io.Reader) (DataBlocks, error) {
	rest := len(bytes) % blockLength
	Blocks := make([][]byte, len(bytes)/blockLength+1)
	for index := 0; index < len(Blocks)-1; index += 1 {
		Blocks[index] = make([]byte, blockLength)
		copy(Blocks[index], bytes[index*blockLength:])
	}
	Blocks[len(Blocks)-1] = make([]byte, blockLength)
	copy(Blocks[len(Blocks)-1], bytes[len(bytes)-rest:])
	_, err := io.ReadFull(reader, Blocks[len(Blocks)-1][rest:])
	return DataBlocks{rest, Blocks}, err
}

func (Blocks *DataBlocks) ToByteArray() []byte {
	blockLength := len(Blocks.Blocks[0])
	bytes := make([]byte, (len(Blocks.Blocks)-1)*blockLength+Blocks.Rest)
	for partIndex, part := range Blocks.Blocks[:len(Blocks.Blocks)-1] {
		copy(bytes[partIndex*len(part):], part)
	}
	copy(bytes[len(bytes)-Blocks.Rest:], Blocks.Blocks[len(Blocks.Blocks)-1])
	return bytes
}
