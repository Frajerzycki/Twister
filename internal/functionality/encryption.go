package functionality

import (
	"crypto/rand"
	"github.com/ikcilrep/gonse/pkg/nse"
	"github.com/ikcilrep/twister/internal/parser"
	"io"
	"math/big"
)

const saltSize int = 16
const blockSize byte = byte(32)

func writeEncryptedBlockWithRecoveryData(writer io.Writer, encryptedBlock []int64, IV []int8) error {
	encryptedBlockBytes := nse.Int64sToBytes(encryptedBlock)
	encryptedBlockLengthBytes := nse.Int64ToBytes(int64(len(encryptedBlockBytes)))

	_, err := writer.Write(encryptedBlockLengthBytes)
	if err != nil {
		return err
	}

	_, err = writer.Write(encryptedBlockBytes)
	if err != nil {
		return err
	}

	_, err = writer.Write(nse.Int8sToBytes(IV))
	if err != nil {
		return err
	}
	return nil
}

func Encrypt(key *big.Int, arguments *parser.Arguments) error {
	salt := make([]byte, saltSize)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}

	bitsToRotate, bytesToRotate, derivedKey, err := nse.DeriveKey(key, salt, int(blockSize))
	if err != nil {
		return err
	}

	_, err = arguments.DataWriter.Write(salt)
	if err != nil {
		return err
	}

	for err == nil {
		var bytesRead int
		block := make([]byte, blockSize)
		bytesRead, err = io.ReadFull(arguments.DataReader, block)
		if err != nil {
			rest, err1 := io.ReadFull(rand.Reader, block[bytesRead:])
			if err1 != nil {
				return err1
			}
			block[len(block)-1] = byte(rest)
		}

		encryptedBlock, IV, err1 := nse.EncryptWithAlreadyDerivedKey(block, derivedKey, bitsToRotate, bytesToRotate)
		if err1 != nil {
			return err1
		}

		err1 = writeEncryptedBlockWithRecoveryData(arguments.DataWriter, encryptedBlock, IV)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func retrieveDataFromReader(reader io.Reader) (encryptedBlock []int64, IV []int8, err error) {
	encryptedBlockBytesLength, _, err := nse.BytesToInt64FromReader(reader)
	if err != nil {
		return nil, nil, err
	}

	encryptedBlockBytes := make([]byte, int(encryptedBlockBytesLength))
	_, err = io.ReadFull(reader, encryptedBlockBytes)
	if err != nil {
		return nil, nil, err
	}

	encryptedBlock, err = nse.BytesToInt64s(encryptedBlockBytes)
	if err != nil {
		return nil, nil, err
	}

	IVBytes := make([]byte, int(blockSize))
	_, err = io.ReadFull(reader, IVBytes)
	if err != nil {
		return nil, nil, err
	}

	IV = nse.BytesToInt8s(IVBytes)

	return
}

func Decrypt(key *big.Int, arguments *parser.Arguments) error {
	salt := make([]byte, saltSize)
	_, err := io.ReadFull(arguments.DataReader, salt)
	if err != nil {
		return err
	}

	bitsToRotate, bytesToRotate, derivedKey, err := nse.DeriveKey(key, salt, int(blockSize))
	if err != nil {
		return err
	}
	var block []byte
	for {
		encryptedBlock, IV, err := retrieveDataFromReader(arguments.DataReader)
		if err != nil {
			rest := int(block[len(block)-1])
			_, err = arguments.DataWriter.Write(block[:len(block)-rest])
			if err != nil {
				return err
			}
			break
		}

		_, err = arguments.DataWriter.Write(block)
		if err != nil {
			return err
		}

		block, err = nse.DecryptWithAlreadyDerivedKey(encryptedBlock, IV, derivedKey, bitsToRotate, bytesToRotate)
		if err != nil {
			return err
		}
	}
	return nil
}
