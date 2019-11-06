package parser

import (
	"os"
	"strconv"
	"strings"
)

func ParseArguments(arguments *Arguments) ([]*os.File, error) {
	var err error
	var hasKeySizeBeenChanged bool
	var hasDataReaderBeenChanged bool
	var hasDataWriterBeenChanged bool
	var hasKeyWriterBeenChanged bool
	var isKeyReadedFromFile bool
	var file *os.File
	files := make([]*os.File, 3)
	filesIndex := 0
	for index := 2; index < len(os.Args); index++ {
		err = nil
		argument := os.Args[index]
		switch argument {
		case "-s":
			if hasKeySizeBeenChanged {
				return nil, &manyParameterValuesError{"Key size"}
			}
			arguments.KeySize, err = parseKeySize(&index)
			hasKeySizeBeenChanged = true
		case "-k":
			if arguments.KeyInput.Reader != nil {
				return nil, &manyParameterValuesError{"Key"}
			}
			arguments.KeyInput.Reader = strings.NewReader(getCommandLineArgument(&index))
		case "-kf":
			if arguments.KeyInput.Reader != nil {
				return nil, &manyParameterValuesError{"Key"}
			}
			file, err = getFileReader(getCommandLineArgument(&index))
			arguments.KeyInput.Reader = file
			files[filesIndex] = file
			filesIndex++
			isKeyReadedFromFile = true
		case "-i":
			if hasDataReaderBeenChanged {
				return nil, &manyParameterValuesError{"Data input"}
			}
			file, err = getFileReader(getCommandLineArgument(&index))
			arguments.DataInput.Reader = file
			files[filesIndex] = file
			filesIndex++
			hasDataReaderBeenChanged = true
		case "-o":
			file, err = getFileWriter(getCommandLineArgument(&index))
			files[filesIndex] = file
			filesIndex++

			if os.Args[1] == "-g" {
				if hasKeyWriterBeenChanged {
					return nil, &manyParameterValuesError{"Key output"}
				}
				arguments.KeyOutput.Writer = file
				hasKeyWriterBeenChanged = true
			} else {
				if hasDataWriterBeenChanged {
					return nil, &manyParameterValuesError{"Data output"}
				}
				arguments.DataOutput.Writer = file
				hasDataWriterBeenChanged = true
			}
		default:
			if submatches := formatArgumentRegexp.FindStringSubmatch(argument); submatches != nil {
				err = parseFormatType(submatches, arguments)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	if !isKeyReadedFromFile && arguments.KeyInput.IsBinary {
		return nil, binaryInputError
	}

	if !hasDataWriterBeenChanged && arguments.DataOutput.IsBinary {
		return nil, binaryOutputError
	}
	return files[:filesIndex], nil
}

func parseKeySize(index *int) (int, error) {
	argument := getCommandLineArgument(index)
	size, err := strconv.Atoi(argument)
	if err == nil {
		return size, nil
	}
	return 0, err
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
