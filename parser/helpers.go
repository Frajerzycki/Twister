package parser

import (
	"os"
)

func getCommandLineArgument(index *int) string {
	(*index)++
	var argument string
	if *index < len(os.Args) {
		argument = os.Args[*index]
	}
	return argument
}

func getFileReader(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getFileWriter(path string) (*os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
