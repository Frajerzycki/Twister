package parser

import (
	"os"
	"strings"
)

func ParseArguments(arguments *Arguments) error {
	var err error
	for index := 2; index < len(os.Args); index++ {
		err = nil
		argument := os.Args[index]
		switch argument {
		case "-s":
			arguments.KeySize, err = parseKeySize(&index)
		case "-k":
			arguments.KeyInput.Reader = strings.NewReader(getCommandLineArgument(&index))
			arguments.KeyInput.IsText = true
		case "-kf":
			arguments.KeyInput.Reader, err = getFileReader(getCommandLineArgument(&index))
		default:
			if submatches := formatArgumentRegexp.FindStringSubmatch(argument); submatches != nil {
				if submatches[2] == "i" {
					arguments.DataInput.IsText = submatches[3] == "d"
					arguments.KeyInput.IsText = submatches[3] == "k"
				} else {
					arguments.DataOutput.IsText = submatches[3] == "d"
					arguments.KeyOutput.IsText = submatches[3] == "k"
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
