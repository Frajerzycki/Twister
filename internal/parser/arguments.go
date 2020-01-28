package parser

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

type Input struct {
	IsBinary bool
	Reader   io.Reader
}

type Output struct {
	IsBinary bool
	Writer   io.Writer
}

type Arguments struct {
	DataInput  *Input
	DataOutput *Output
	KeyInput   *Input
	KeyOutput  *Output
	KeySize    int
}

func (arguments *Arguments) GetKey() (*big.Int, error) {
	key := new(big.Int)
	keyBytes, err := ioutil.ReadAll(arguments.KeyInput.Reader)

	switch {
	case err != nil:
		return nil, err
	case arguments.KeyInput.IsBinary:
		key.SetBytes(keyBytes)
	default:
		lastIndex := len(keyBytes) - 1
		if keyBytes[lastIndex] == '\n' {
			keyBytes = keyBytes[:lastIndex]
		}
		decodedBytes, err := base64.StdEncoding.DecodeString(string(keyBytes))
		if err != nil {
			return nil, err
		}
		key.SetBytes(decodedBytes)
	}
	return key, nil
}

func (arguments *Arguments) WriteToDataOutput(data []byte) {
	if arguments.DataOutput.IsBinary {
		arguments.DataOutput.Writer.Write(data)
	} else {
		io.WriteString(arguments.DataOutput.Writer, base64.StdEncoding.EncodeToString(data))
	}
}

func (arguments *Arguments) WriteToNonBinaryDataOutput(data string) {
	if !arguments.DataOutput.IsBinary {
		io.WriteString(arguments.DataOutput.Writer, data)
	}
}

func NewArguments() *Arguments {
	return &Arguments{DataInput: &Input{Reader: os.Stdin}, DataOutput: &Output{Writer: os.Stdout}, KeyInput: &Input{}, KeyOutput: &Output{Writer: os.Stdout}, KeySize: 32}
}
