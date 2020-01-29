package parser

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"math/big"
	"os"
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

func NewArguments() *Arguments {
	return &Arguments{DataReader: os.Stdin, DataWriter: os.Stdout, KeyWriter: os.Stdout, KeySize: 32}
}
