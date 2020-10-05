package functionality

import (
	"crypto/rand"
	"io"
	"math/big"

	"github.com/ikcilrep/gonse/pkg/nse"
	"github.com/ikcilrep/twister/internal/parser"
)

const saltSize int = 16
const blockSize int = 256

func writeEncryptedBlockWithRecoveryData(writer io.Writer, encryptedBlock []int64, IV []int8, salt []byte) (int, error) {
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

	bytesWritten4, err := writer.Write(salt)
	if err != nil {
		return bytesWritten1 + bytesWritten2 + bytesWritten3 + bytesWritten4, err
	}

	return bytesWritten1 + bytesWritten2 + bytesWritten3 + bytesWritten4, nil
}

func Encrypt(key *big.Int, arguments *parser.Arguments) (bytesRead int64, bytesWritten int64, err error) {

	bytesRead = int64(0)
	bytesWritten = int64(0)
	shouldContinue := true
	for shouldContinue {
		var bytesRead1, bytesWritten1 int
		bytesRead1, bytesWritten1, shouldContinue, err = encryptBlock(arguments, key)
		bytesRead += int64(bytesRead1)
		bytesWritten += int64(bytesWritten1)
		if err != nil {
			return bytesRead, bytesWritten, err
		}
	}
	err = arguments.DataWriter.Close()
	return
}

func pad(block []byte, filledBytes, blockSize int) error {
	var rest int
	rest, err := io.ReadFull(rand.Reader, block[filledBytes:])
	if err != nil {
		return err
	}
	block[len(block)-1] = byte(rest)
	return nil
}

func encryptBlock(arguments *parser.Arguments, key *big.Int) (bytesRead int, bytesWritten int, shouldContinue bool, err error) {
	salt := make([]byte, saltSize)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		return 0, 0, false, err
	}
	derivedKey, err := nse.DeriveKey(key, salt, blockSize)
	if err != nil {
		return 0, 0, false, err
	}

	block := make([]byte, blockSize)
	bytesRead, err = io.ReadFull(arguments.DataReader, block)
	shouldContinue = true
	if err != nil {
		shouldContinue = false
		err = pad(block, bytesRead, blockSize)
		if err != nil {
			return
		}
	}

	encryptedBlock, IV, err := nse.Encrypt(block, derivedKey)
	if err != nil {
		return
	}

	bytesWritten, err = writeEncryptedBlockWithRecoveryData(arguments.DataWriter, encryptedBlock, IV, salt)
	if err != nil {
		return
	}

	return
}
