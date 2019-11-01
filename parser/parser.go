package parser

import (
	"os"
)

func ParseArguments(parameters *Arguments) error {
	var err error
	for index := 2; index < len(os.Args); index++ {
		err = nil
		argument := os.Args[index]
		switch argument {
		case "-s":
			parameters.KeySize, err = parseKeySize(&index)
		case "-i":
			parameters.Data, err = getDataFromStd()
		case "-k":
			parameters.Key, err = parseKey(&index)
		case "-kf":
			parameters.Key, err = parseKeyFromFile(&index)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
