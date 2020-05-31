package parser

import (
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
)

// Arguments represents arguments of Twister.
type Arguments struct {
	DataWriter io.WriteCloser
	DataReader io.Reader
	KeyWriter  io.Writer
	KeyReader  io.Reader
	KeySize    int
}

// GetKey reads and returns secret key.
func (arguments *Arguments) GetKey() (*big.Int, error) {
	key := new(big.Int)
	keyBytes, err := ioutil.ReadAll(arguments.KeyReader)

	if err != nil {
		return nil, err
	}

	lastIndex := len(keyBytes) - 1
	if keyBytes[lastIndex] == '\n' {
		keyBytes = keyBytes[:lastIndex]
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(string(keyBytes))
	if err != nil {
		return nil, err
	}
	key.SetBytes(decodedBytes)
	return key, nil
}

// VerifyForNSETransformation verifies whether correct arguments are used for encryption or decryption.
func (arguments *Arguments) VerifyForNSETransformation() error {
	if arguments.DataWriter == nil {
		return errors.New("Output is not set.")
	}
	if arguments.DataReader == nil {
		return errors.New("Input is not set.")
	}

	if arguments.KeyReader == nil {
		return errors.New("Key is not set.")
	}
	return nil
}

// VerifyForKeyGeneration verifies whether correct arguments are used for key generation.
func (arguments *Arguments) VerifyForKeyGeneration() error {
	if arguments.KeyWriter == nil {
		return errors.New("Output is not set.")
	}
	return nil
}

// NewArguments creates new instance of arguments with 32 bytes key.
func NewArguments() *Arguments {
	return &Arguments{KeySize: 32}
}
