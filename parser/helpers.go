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
			if arguments.DataInput.IsText {
				return &manyParameterValuesError{"Is input data text"}
			}
			arguments.DataInput.IsText = true
		} else {
			if arguments.KeyInput.IsText {
				return &manyParameterValuesError{"Is input key text"}
			}
			arguments.KeyInput.IsText = true
		}
	} else {
		if submatches[3] == "d" {
			if arguments.DataOutput.IsText {
				return &manyParameterValuesError{"Is output data text"}
			}
			arguments.DataOutput.IsText = true
		} else {
			if arguments.KeyOutput.IsText {
				return &manyParameterValuesError{"Is output key text"}
			}
			arguments.KeyOutput.IsText = true
		}
	}
	return nil
}
