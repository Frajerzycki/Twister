package parser

import (
	"io"
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
	KeySize    uint
}

func NewArguments() *Arguments {
	return &Arguments{DataInput: &Input{Reader: os.Stdin}, DataOutput: &Output{Writer: os.Stdout}, KeyInput: &Input{}, KeyOutput: &Output{Writer: os.Stdout}, KeySize: uint(256)}
}
