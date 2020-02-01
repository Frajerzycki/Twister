package functionality

import (
	"crypto/rand"
	"github.com/ikcilrep/gonse/pkg/nse"
	"github.com/ikcilrep/twister/internal/parser"
	"io"
	"math/big"
)

const saltSize int = 16
const blockSize int = 256

func writeEncryptedBlockWithRecoveryData(writer io.Writer, encryptedBlock []int64, IV []int8) (int, error) {
	encryptedBlockBytes := nse.Int64sToBytes(encryptedBlock)
	encryptedBlockLengthBytes := nse.Int64ToBytes(int64(len(encryptedBlockBytes)))

	bytesWritten1, err := writer.Write(encryptedBlockLengthBytes)
	if err != nil {
		return bytesWritten1, err
	}

	bytesWritten2, err := writer.Write(encryptedBlockBytes)
	if err != nil {
		return bytesWritten1 + bytesWritten2, err
	}

	bytesWritten3, err := writer.Write(nse.Int8sToBytes(IV))
	if err != nil {
		return bytesWritten1 + bytesWritten2 + bytesWritten3, err
	}
	return bytesWritten1 + bytesWritten2 + bytesWritten3, nil
}

func Encrypt(key *big.Int, arguments *parser.Arguments) (bytesRead int64, bytesWritten int64, err error) {
	salt := make([]byte, saltSize)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		return int64(0), int64(0), err
	}

	derivedKey, err := nse.DeriveKey(key, salt, blockSize)
	if err != nil {
		return int64(0), int64(0), err
	}

	saltBytesWritten, err := arguments.DataWriter.Write(salt)
	if err != nil {
		return int64(0), int64(saltBytesWritten), err
	}
	bytesRead = int64(0)
	bytesWritten = int64(saltSize)
	for err == nil {
		var bytesRead1 int
		block := make([]byte, blockSize)
		bytesRead1, err = io.ReadFull(arguments.DataReader, block)
		bytesRead += int64(bytesRead1)
		if err != nil {
			rest, err1 := io.ReadFull(rand.Reader, block[bytesRead1:])
			if err1 != nil {
				return bytesRead, bytesWritten, err1
			}
			block[len(block)-1] = byte(rest)
		}

		encryptedBlock, IV, err1 := nse.Encrypt(block, derivedKey)
		if err1 != nil {
			return bytesRead, bytesWritten, err1
		}

		bytesWritten1, err1 := writeEncryptedBlockWithRecoveryData(arguments.DataWriter, encryptedBlock, IV)
		bytesWritten += int64(bytesWritten1)
		if err1 != nil {
			return bytesRead, bytesWritten, err1
		}

	}
	err = nil
	return
}

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
	bytesRead = int64(saltSize)
	bytesWritten = int64(0)
	for {
		encryptedBlock, IV, bytesRead1, err := retrieveDataFromReader(arguments.DataReader)
		bytesRead += int64(bytesRead1)
		if err != nil {
			rest := int(block[len(block)-1])
			_, err = arguments.DataWriter.Write(block[:len(block)-rest])
			if err != nil {
				return bytesRead, bytesWritten, err
			}
			break
		}

		bytesWritten1, err := arguments.DataWriter.Write(block)
		bytesWritten += int64(bytesWritten1)
		if err != nil {
			return bytesRead, bytesWritten, err
		}

		block, err = nse.Decrypt(encryptedBlock, IV, derivedKey)
		if err != nil {
			return bytesRead, bytesWritten, err
		}
	}
	err = nil
	return
}
