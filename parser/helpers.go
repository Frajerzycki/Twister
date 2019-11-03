package parser

import (
	"io"
	"os"
	"strconv"
)

func getCommandLineArgument(index *int) string {
	(*index)++
	var argument string
	if *index < len(os.Args) {
		argument = os.Args[*index]
	}
	return argument
}

func parseKeySize(index *int) (uint, error) {
	argument := getCommandLineArgument(index)
	size, err := strconv.Atoi(argument)
	if err == nil {
		return uint(size), nil
	}
	return 0, err
}

func getFileReader(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
