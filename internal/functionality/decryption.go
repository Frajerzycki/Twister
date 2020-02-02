package functionality

import (
	"github.com/ikcilrep/gonse/pkg/nse"
	"github.com/ikcilrep/twister/internal/parser"
	"io"
	"math/big"
)

func retrieveDataFromReader(reader io.Reader) (encryptedBlock []int64, IV []int8, bytesRead int, err error) {
	encryptedBlockBytesLength, bytesRead1, err := nse.BytesToInt64FromReader(reader)
	bytesRead += bytesRead1
	if err != nil {
		return nil, nil, bytesRead, err
	}

	encryptedBlockBytes := make([]byte, int(encryptedBlockBytesLength))
	bytesRead1, err = io.ReadFull(reader, encryptedBlockBytes)
	bytesRead += bytesRead1
	if err != nil {
		return nil, nil, bytesRead, err
	}

	encryptedBlock, err = nse.BytesToInt64s(encryptedBlockBytes)
	if err != nil {
		return nil, nil, bytesRead, err
	}

	IVBytes := make([]byte, blockSize)
	bytesRead1, err = io.ReadFull(reader, IVBytes)
	bytesRead += bytesRead1
	if err != nil {
		return nil, nil, bytesRead, err
	}

	IV = nse.BytesToInt8s(IVBytes)
	return
}

func unpad(block []byte, blockSize int) []byte {
	rest := int(block[len(block)-1])
	return block[:len(block)-rest]
}

func decryptBlock(arguments *parser.Arguments, derivedKey *nse.NSEKey, lastBlock []byte) (bytesRead int, bytesWritten int, block []byte, shouldContinue bool, err error) {
	var bytesWritten1 int
	encryptedBlock, IV, bytesRead, err := retrieveDataFromReader(arguments.DataReader)

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

	block, err = nse.Decrypt(encryptedBlock, IV, derivedKey)
	if err != nil {
		return
	}

	return
}

func Decrypt(key *big.Int, arguments *parser.Arguments) (bytesRead int64, bytesWritten int64, err error) {
	salt := make([]byte, saltSize)
	saltBytesRead, err := io.ReadFull(arguments.DataReader, salt)
	if err != nil {
		return int64(saltBytesRead), int64(0), err
	}

	derivedKey, err := nse.DeriveKey(key, salt, blockSize)
	if err != nil {
		return int64(saltSize), int64(0), err
	}
	var block []byte
	shouldContinue := true
	bytesRead = int64(saltSize)
	bytesWritten = int64(0)

	for shouldContinue {
		var bytesRead1, bytesWritten1 int

		bytesRead1, bytesWritten1, block, shouldContinue, err = decryptBlock(arguments, derivedKey, block)
		bytesRead += int64(bytesRead1)
		bytesWritten += int64(bytesWritten1)
		if err != nil {
			return bytesRead, bytesWritten, err
		}
	}
	return
}
