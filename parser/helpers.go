package parser

import (
	"io"
	"os"
	"strconv"
)

func getCommandLineArgument(index *int) string {
	(*index)++
	var argument string
	if *index < len(os.Args) {
		argument = os.Args[*index]
	}
	return argument
}

func parseKeySize(index *int) (uint, error) {
	argument := getCommandLineArgument(index)
	size, err := strconv.Atoi(argument)
	if err == nil {
		return uint(size), nil
	}
	return 0, err
}

func getFileReader(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func parseFormatType(submatches []string, arguments *Arguments) error {
	if submatches[2] == "i" {
		if submatches[3] == "d" {
			if arguments.DataInput.IsBinary {
				return &manyParameterValuesError{"Is input data binary"}
			}
			arguments.DataInput.IsBinary = true
		} else {
			if arguments.KeyInput.IsBinary {
				return &manyParameterValuesError{"Is input key binary"}
			}
			arguments.KeyInput.IsBinary = true
		}
	} else {
		if submatches[3] == "d" {
			if arguments.DataOutput.IsBinary {
				return &manyParameterValuesError{"Is output data binary"}
			}
			arguments.DataOutput.IsBinary = true
		} else {
			if arguments.KeyOutput.IsBinary {
				return &manyParameterValuesError{"Is output key binary"}
			}
			arguments.KeyOutput.IsBinary = true
		}
	}
	return nil
}
