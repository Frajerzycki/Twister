package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Frajerzycki/GONSE"
	"github.com/Frajerzycki/Twister/parser"
	"math/big"
)

const saltSize int = 16

var wrongCiphertextFormatError error = errors.New("Wrong format of ciphertext given to encrypt.")

func encrypt(data []byte, key *big.Int, arguments *parser.Arguments) error {
	IV, err := nse.GenerateIV(len(data))
	if err != nil {
		return err
	}
	var salt []byte
	salt, err = randomBytes(saltSize)
	if err != nil {
		return err
	}

	ciphertext, err := nse.Encrypt(data, salt, IV, key)
	if err != nil {
		return err
	}
	bytes, _ := nse.Int64sToBytes(ciphertext)
	encryptedToBytesLength := len(bytes)
	bytes = append(bytes, nse.Int8sToBytes(IV)...)
	buffer := make([]byte, 8)
	binary.PutUvarint(buffer, uint64(encryptedToBytesLength))
	bytes = append(bytes, buffer...)
	bytes = append(bytes, salt...)

	if arguments.DataOutput.IsBinary {
		arguments.DataOutput.Writer.Write(bytes)
	} else {
		arguments.DataOutput.Writer.Write([]byte(fmt.Sprintf("%v\n", base64.StdEncoding.EncodeToString(bytes))))
	}
	return nil
}

func decrypt(data []byte, key *big.Int, arguments *parser.Arguments) error {
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

func generateKey(keySize int, arguments *parser.Arguments) error {
	keyBytes, err := randomBytes(keySize)

	if err != nil {
		return err
	}

	key := new(big.Int)
	key.SetBytes(keyBytes)

	if arguments.KeyOutput.IsBinary {
		arguments.KeyOutput.Writer.Write(key.Bytes())
	} else {
		arguments.KeyOutput.Writer.Write([]byte(fmt.Sprintf("%v\n", key.Text(parser.KeyBase))))
	}
	return nil
}
