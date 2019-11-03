package parser

import (
	"os"
	"strings"
)

func ParseArguments(arguments *Arguments) error {
	var err error
	var hasKeySizeBeenChanged bool
	var isKeyReadedFromFile bool
	for index := 2; index < len(os.Args); index++ {
		err = nil
		argument := os.Args[index]
		switch argument {
		case "-s":
			if hasKeySizeBeenChanged {
				return &manyParameterValuesError{"Key size"}
			}
			arguments.KeySize, err = parseKeySize(&index)
			hasKeySizeBeenChanged = true
		case "-k":
			if arguments.KeyInput.Reader != nil {
				return &manyParameterValuesError{"Key"}
			}
			arguments.KeyInput.Reader = strings.NewReader(getCommandLineArgument(&index))
		case "-kf":
			if arguments.KeyInput.Reader != nil {
				return &manyParameterValuesError{"Key"}
			}
			arguments.KeyInput.Reader, err = getFileReader(getCommandLineArgument(&index))
			isKeyReadedFromFile = true
		default:
			if submatches := formatArgumentRegexp.FindStringSubmatch(argument); submatches != nil {
				err = parseFormatType(submatches, arguments)
			}
		}
		if err != nil {
			return err
		}
	}
	if !isKeyReadedFromFile && arguments.KeyInput.IsBinary {
		return binaryInputError
	}

	return nil
}
