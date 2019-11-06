package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/Frajerzycki/GONSE"
	"github.com/Frajerzycki/Twister/parser"
	"math/big"
)

const saltSize int = 16

func encrypt(data []byte, key *big.Int, arguments *parser.Arguments) error {
	IV, err := nse.GenerateIV(len(data))
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
	bytes = append(bytes, nse.Int8sToBytes(IV)...)
	buffer := make([]byte, 8)
	binary.PutUvarint(buffer, uint64(len(data)))
	bytes = append(bytes, buffer...)
	bytes = append(bytes, salt...)
	// Only for testing
	if arguments.DataOutput.IsBinary {
		arguments.DataOutput.Writer.Write(bytes)
	} else {
		arguments.DataOutput.Writer.Write([]byte(fmt.Sprintf("%v\n", base64.StdEncoding.EncodeToString(bytes))))
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
