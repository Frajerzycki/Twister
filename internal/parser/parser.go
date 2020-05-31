package parser

import (
	"github.com/ikcilrep/twister/internal/files"
	"os"
	"strconv"
)

const chunkSize int = 1048576

// ParseArguments parses arguments into Arguments struct.
func ParseArguments(arguments *Arguments) ([]*os.File, error) {
	var err error
	var hasKeySizeBeenChanged bool
	var hasDataReaderBeenChanged bool
	var hasDataWriterBeenChanged bool
	var hasKeyWriterBeenChanged bool
	var file *os.File
	openedFiles := make([]*os.File, 3)
	openedFilesIndex := 0
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
			file, err = getFileReader(getCommandLineArgument(&index))
			arguments.KeyReader = file
			openedFiles[openedFilesIndex] = file
			openedFilesIndex++
		case "-i":
			if hasDataReaderBeenChanged {
				return nil, &manyParameterValuesError{"Data input"}
			}
			file, err = getFileReader(getCommandLineArgument(&index))
			arguments.DataReader, err = files.NewChunkedReader(file, chunkSize)
			if err != nil {
				return nil, err
			}
			openedFiles[openedFilesIndex] = file
			openedFilesIndex++
			hasDataReaderBeenChanged = true
		case "-o":
			file, err = getFileWriter(getCommandLineArgument(&index))
			openedFiles[openedFilesIndex] = file
			openedFilesIndex++

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
				arguments.DataWriter, err = files.NewChunkedWriter(file, chunkSize)
				if err != nil {
					return nil, err
				}
				hasDataWriterBeenChanged = true
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return openedFiles[:openedFilesIndex], nil
}

func parseKeySize(index *int) (int, error) {
	argument := getCommandLineArgument(index)
	size, err := strconv.Atoi(argument)
	if err == nil {
		return size, nil
	}
	return 0, err
}
