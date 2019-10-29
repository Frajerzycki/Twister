package main

import (
	"math/big"
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
	return 0, &intFormatError{os.Args[*index], keyBase}
}

func parseKey(index *int) (*big.Int, error) {
	(*index)++
	key := new(big.Int)
	var ok bool
	if *index < len(os.Args) {
		if key, ok = key.SetString(os.Args[*index], keyBase); !ok {
			return nil, &intFormatError{os.Args[*index], keyBase}
		}
	}
	return key, nil
}
