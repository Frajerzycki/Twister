package parser

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

var formatArgumentRegexp *regexp.Regexp = regexp.MustCompile("(-b)([io])([kd])")

func ParseArguments(arguments *Arguments) ([]*os.File, error) {
	var err error
	var hasKeySizeBeenChanged bool
	var hasDataReaderBeenChanged bool
	var hasDataWriterBeenChanged bool
	var hasKeyWriterBeenChanged bool
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
			if arguments.KeyReader != nil {
				return nil, &manyParameterValuesError{"Key"}
			}
			arguments.KeyReader = strings.NewReader(getCommandLineArgument(&index))
		case "-kf":
			if arguments.KeyReader != nil {
				return nil, &manyParameterValuesError{"Key"}
			}
			file, err = getFileReader(getCommandLineArgument(&index))
			arguments.KeyReader = file
			files[filesIndex] = file
			filesIndex++
		case "-i":
			if hasDataReaderBeenChanged {
				return nil, &manyParameterValuesError{"Data input"}
			}
			file, err = getFileReader(getCommandLineArgument(&index))
			arguments.DataReader = file
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
				arguments.KeyWriter = file
				hasKeyWriterBeenChanged = true
			} else {
				if hasDataWriterBeenChanged {
					return nil, &manyParameterValuesError{"Data output"}
				}
				arguments.DataWriter = file
				hasDataWriterBeenChanged = true
			}
		}
		if err != nil {
			return nil, err
		}
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
