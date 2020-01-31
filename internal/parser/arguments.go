package parser

import (
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
)

type Arguments struct {
	DataWriter io.Writer
	DataReader io.Reader
	KeyWriter  io.Writer
	KeyReader  io.Reader
	KeySize    int
}

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

func (arguments *Arguments) VerifyForKeyGeneration() error {
	if arguments.KeyWriter == nil {
		return errors.New("Output is not set.")
	}
	return nil
}

func NewArguments() *Arguments {
	return &Arguments{KeySize: 32}
}
