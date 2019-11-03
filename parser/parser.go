package parser

import (
	"os"
	"strings"
)

func ParseArguments(parameters *Arguments) error {
	var err error
	for index := 2; index < len(os.Args); index++ {
		err = nil
		argument := os.Args[index]
		switch argument {
		case "-s":
			parameters.KeySize, err = parseKeySize(&index)
		case "-k":
			parameters.KeyInput.Reader = strings.NewReader(getCommandLineArgument(&index))
			parameters.KeyInput.IsText = true
		case "-kf":
			parameters.KeyInput.Reader, err = getFileReader(getCommandLineArgument(&index))
		default:
			if submatches := formatArgumentRegexp.FindStringSubmatch(argument); submatches != nil {
				if submatches[2] == "i" {
					parameters.DataInput.IsText = submatches[3] == "d"
					parameters.KeyInput.IsText = submatches[3] == "k"
				} else {
					parameters.DataOutput.IsText = submatches[3] == "d"
					parameters.KeyOutput.IsText = submatches[3] == "k"
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
