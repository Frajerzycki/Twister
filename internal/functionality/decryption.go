package functionality

import (
	"io"
	"math/big"

	"github.com/ikcilrep/gonse/pkg/nse"
	"github.com/ikcilrep/twister/internal/parser"
)

func retrieveDataFromReader(reader io.Reader) (encryptedBlock []int64, IV []int8, salt []byte, bytesRead int, err error) {
	encryptedBlockBytesLength, bytesRead1, err := nse.BytesToInt64FromReader(reader)
	bytesRead += bytesRead1
	if err != nil {
		return nil, nil, nil, bytesRead, err
	}

	encryptedBlockBytes := make([]byte, int(encryptedBlockBytesLength))
	bytesRead1, err = io.ReadFull(reader, encryptedBlockBytes)
	bytesRead += bytesRead1
	if err != nil {
		return nil, nil, nil, bytesRead, err
	}

	encryptedBlock, err = nse.BytesToInt64s(encryptedBlockBytes)
	if err != nil {
		return nil, nil, nil, bytesRead, err
	}

	IVBytes := make([]byte, blockSize)
	bytesRead1, err = io.ReadFull(reader, IVBytes)
	bytesRead += bytesRead1
	if err != nil {
		return nil, nil, nil, bytesRead, err
	}

	IV = nse.BytesToInt8s(IVBytes)

	salt = make([]byte, saltSize)
	bytesRead1, err = io.ReadFull(reader, salt)
	bytesRead += bytesRead1

	return
}

func unpad(block []byte, blockSize int) []byte {
	rest := int(block[len(block)-1])
	return block[:len(block)-rest]
}

func decryptBlock(arguments *parser.Arguments, key *big.Int, lastBlock []byte) (bytesRead int, bytesWritten int, block []byte, shouldContinue bool, err error) {
	var bytesWritten1 int
	encryptedBlock, IV, salt, bytesRead, err := retrieveDataFromReader(arguments.DataReader)

	shouldContinue = true
	if err != nil {
		shouldContinue = false
		unpaddedLastBlock := unpad(lastBlock, blockSize)
		bytesWritten1, err = arguments.DataWriter.Write(unpaddedLastBlock)
		bytesWritten += bytesWritten1
		if err != nil {
			return
		}
		return
	}

	bytesWritten1, err = arguments.DataWriter.Write(lastBlock)
	bytesWritten += bytesWritten1
	if err != nil {
		return
	}

	derivedKey, err := nse.DeriveKey(key, salt, blockSize)
	if err != nil {
		return bytesRead, 0, nil, false, err
	}

	block, err = nse.Decrypt(encryptedBlock, IV, derivedKey)
	if err != nil {
		return
	}

	return
}

// Decrypt decrypts file with key and writes plaintext to another file.
func Decrypt(key *big.Int, arguments *parser.Arguments) (bytesRead int64, bytesWritten int64, err error) {
	var block []byte
	shouldContinue := true
	bytesRead = int64(0)
	bytesWritten = int64(0)

	for shouldContinue {
		var bytesRead1, bytesWritten1 int

		bytesRead1, bytesWritten1, block, shouldContinue, err = decryptBlock(arguments, key, block)
		bytesRead += int64(bytesRead1)
		bytesWritten += int64(bytesWritten1)
		if err != nil {
			return bytesRead, bytesWritten, err
		}
	}
	err = arguments.DataWriter.Close()
	return
}
