package parser

import (
	"os"
	"strings"
)

func ParseArguments(arguments *Arguments) ([]*os.File, error) {
	var err error
	var hasKeySizeBeenChanged bool
	var hasDataReaderBeenChanged bool
	var hasDataWriterBeenChanged bool
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
			if hasDataWriterBeenChanged {
				return nil, &manyParameterValuesError{"Data output"}
			}
			file, err = getFileWriter(getCommandLineArgument(&index))
			arguments.DataOutput.Writer = file
			files[filesIndex] = file
			filesIndex++
			hasDataWriterBeenChanged = true
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
