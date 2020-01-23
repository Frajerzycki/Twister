package functionality

import "io"

type DataBlocks struct {
	Rest   byte
	Blocks [][]byte
}

func ByteArrayToDataBlocks(bytes []byte, blockLength byte, reader io.Reader) (DataBlocks, error) {
	rest := len(bytes) % int(blockLength)
	Blocks := make([][]byte, len(bytes)/int(blockLength)+1)
	for index := 0; index < len(Blocks)-1; index += 1 {
		Blocks[index] = make([]byte, int(blockLength))
		copy(Blocks[index], bytes[index*int(blockLength):])
	}
	Blocks[len(Blocks)-1] = make([]byte, int(blockLength))
	copy(Blocks[len(Blocks)-1], bytes[len(bytes)-rest:])
	_, err := io.ReadFull(reader, Blocks[len(Blocks)-1][rest:])
	return DataBlocks{byte(rest), Blocks}, err
}

func (Blocks *DataBlocks) ToByteArray() []byte {
	blockLength := len(Blocks.Blocks[0])
	bytes := make([]byte, (len(Blocks.Blocks)-1)*blockLength+int(Blocks.Rest))
	for partIndex, part := range Blocks.Blocks[:len(Blocks.Blocks)-1] {
		copy(bytes[partIndex*len(part):], part)
	}
	copy(bytes[len(bytes)-int(Blocks.Rest):], Blocks.Blocks[len(Blocks.Blocks)-1])
	return bytes
}
