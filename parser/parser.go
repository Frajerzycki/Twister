package parser

import "os"

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
		default:
			if submatches := formatArgumentRegexp.FindStringSubmatch(argument); submatches != nil {
				if submatches[2] == "i" {
					parameters.IsInputDataText = submatches[3] == "d"
					parameters.IsInputKeyText = submatches[3] == "k"
				} else {
					parameters.IsOutputDataText = submatches[3] == "d"
					parameters.IsOutputKeyText = submatches[3] == "k"
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
