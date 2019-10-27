package main

import (
	"os"
	"strconv"
)

func parseKeySize(index *int) (uint, error) {
	(*index)++
	if *index < len(os.Args) {
		if size, err := strconv.Atoi(os.Args[*index]); err == nil {
			return uint(size), nil
		}
	}
	return 0, &keySizeFormatError{os.Args[*index]}
}
