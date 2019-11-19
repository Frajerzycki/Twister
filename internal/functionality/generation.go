package functionality

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ikcilrep/twister/internal/parser"
	"math/big"
)

var notPositiveLengthError error = errors.New("Length of data which exists has to be positive.")

func randomBytes(length int) ([]byte, error) {
	if length < 0 {
		return nil, notPositiveLengthError
	}
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateKey(arguments *parser.Arguments) error {
	keyBytes, err := randomBytes(arguments.KeySize)

	if err != nil {
		return err
	}

	key := new(big.Int)
	key.SetBytes(keyBytes)

	if arguments.KeyOutput.IsBinary {
		arguments.KeyOutput.Writer.Write(key.Bytes())
	} else {
		arguments.KeyOutput.Writer.Write([]byte(fmt.Sprintf("%v\n", base64.StdEncoding.EncodeToString(key.Bytes()))))
	}
	return nil
}
