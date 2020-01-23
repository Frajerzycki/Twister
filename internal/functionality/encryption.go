package functionality

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ikcilrep/gonse/pkg/nse"
	"github.com/ikcilrep/twister/internal/parser"
	"math/big"
)

const saltSize int = 16
const blockLength int = 32

var wrongCiphertextFormatError error = errors.New("Wrong format of ciphertext given to encrypt.")

func addRecoveryInfoToCiphertext(destiny *[]byte, ciphertext []int64, IV []int8) {
	ciphertextBytes := nse.Int64sToBytes(ciphertext)
	ciphertextLengthBytes := nse.Int64ToBytes(int64(len(ciphertextBytes)))

	*destiny = make([]byte, len(ciphertextBytes)+len(ciphertextLengthBytes)+len(IV))
	copy(*destiny, ciphertextBytes)
	copy((*destiny)[len(ciphertextBytes):], ciphertextLengthBytes)
	copy((*destiny)[len(ciphertextBytes)+len(ciphertextLengthBytes):], nse.Int8sToBytes(IV))

}

func Encrypt(data []byte, key *big.Int, arguments *parser.Arguments) error {
	var salt []byte
	salt, err := randomBytes(saltSize)
	if err != nil {
		return err
	}

	bitsToRotate, bytesToRotate, derivedKey, err := nse.DeriveKey(key, salt, len(data))
	if err != nil {
		return err
	}

	dataBlocks, err := ByteArrayToDataBlocks(data, blockLength, rand.Reader)
	if err != nil {
		return err
	}
	encryptedDataBlocks := DataBlocks{Blocks: make([][]byte, len(dataBlocks.Blocks)), Rest: 0}

	for index, block := range dataBlocks.Blocks {
		ciphertext, IV, err := nse.EncryptWithAlreadyDerivedKey(block, salt, derivedKey, bitsToRotate, bytesToRotate)
		if err != nil {
			return err
		}
		addRecoveryInfoToCiphertext(&encryptedDataBlocks.Blocks[index], ciphertext, IV)
	}
	encryptedBytes := encryptedDataBlocks.ToByteArray()
	encryptedBytes = append(encryptedBytes, salt...)

	if arguments.DataOutput.IsBinary {
		arguments.DataOutput.Writer.Write(encryptedBytes)
	} else {
		arguments.DataOutput.Writer.Write([]byte(fmt.Sprintf("%v\n", base64.StdEncoding.EncodeToString(encryptedBytes))))
	}
	return nil
}

func Decrypt(data []byte, key *big.Int, arguments *parser.Arguments) error {
	var ciphertext []byte
	var err error
	if arguments.DataInput.IsBinary {
		ciphertext = data
	} else {
		ciphertext, err = base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			return err
		}
	}
	salt := ciphertext[len(ciphertext)-saltSize:]
	toDecryptLength, _ := binary.Uvarint(ciphertext[len(ciphertext)-saltSize-8 : len(ciphertext)-saltSize])
	if toDecryptLength > uint64(len(ciphertext)) || toDecryptLength >= uint64(len(ciphertext)-saltSize-8) {
		return wrongCiphertextFormatError
	}
	toDecrypt, err := nse.BytesToInt64s(ciphertext[:toDecryptLength])
	if err != nil {
		return err
	}
	IV := nse.BytesToInt8s(ciphertext[toDecryptLength : len(ciphertext)-saltSize-8])
	decrypted, err := nse.Decrypt(toDecrypt, salt, IV, key)
	switch {
	case err != nil:
		return err
	case arguments.DataOutput.IsBinary:
		arguments.DataOutput.Writer.Write(decrypted)
	default:
		arguments.DataOutput.Writer.Write([]byte(string(decrypted)))
	}
	return nil

}
