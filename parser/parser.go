package parser

import (
	"os"
	"strings"
)

func ParseArguments(arguments *Arguments) error {
	var err error
	var hasKeySizeBeenSet bool
	for index := 2; index < len(os.Args); index++ {
		err = nil
		argument := os.Args[index]
		switch argument {
		case "-s":
			if hasKeySizeBeenSet {
				return &manyParameterValuesError{"Key size"}
			}
			arguments.KeySize, err = parseKeySize(&index)
			hasKeySizeBeenSet = true
		case "-k":
			if arguments.KeyInput.Reader != nil {
				return &manyParameterValuesError{"Key"}
			}
			arguments.KeyInput.Reader = strings.NewReader(getCommandLineArgument(&index))
			arguments.KeyInput.IsText = true
		case "-kf":
			if arguments.KeyInput.Reader != nil {
				return &manyParameterValuesError{"Key"}
			}
			arguments.KeyInput.Reader, err = getFileReader(getCommandLineArgument(&index))
		default:
			if submatches := formatArgumentRegexp.FindStringSubmatch(argument); submatches != nil {
				err = parseFormatType(submatches, arguments)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
