package functionality

import (
	"crypto/rand"
	"errors"
	"github.com/ikcilrep/gonse/pkg/nse"
	"github.com/ikcilrep/twister/internal/parser"
	"io"
	"math/big"
)

const saltSize int = 16
const blockSize byte = byte(32)

var wrongCiphertextFormatError error = errors.New("Wrong format of ciphertext given to encrypt.")

// EncryptedBlock = cipherTextLength + ciphertext + IV
func writeCiphertextBlockWithRecoveryData(arguments *parser.Arguments, ciphertext []int64, IV []int8) {
	ciphertextBytes := nse.Int64sToBytes(ciphertext)
	ciphertextLengthBytes := nse.Int64ToBytes(int64(len(ciphertextBytes)))
	arguments.WriteToDataOutput(ciphertextLengthBytes)
	arguments.WriteToDataOutput(ciphertextBytes)
	arguments.WriteToDataOutput(nse.Int8sToBytes(IV))
}

func Encrypt(key *big.Int, arguments *parser.Arguments) error {
	salt := make([]byte, saltSize)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}

	bitsToRotate, bytesToRotate, derivedKey, err := nse.DeriveKey(key, salt, int(blockSize))

	for err != nil {
		var bytesRead int
		block := make([]byte, blockSize)
		bytesRead, err = io.ReadFull(arguments.DataInput.Reader, block)
		if err != nil {
			_, err := io.ReadFull(rand.Reader, block[bytesRead:])
			if err != nil {
				return err
			}
		}

		ciphertext, IV, err := nse.EncryptWithAlreadyDerivedKey(block, derivedKey, bitsToRotate, bytesToRotate)
		if err != nil {
			return err
		}
		writeCiphertextBlockWithRecoveryData(arguments, ciphertext, IV)
	}
	return nil
}
