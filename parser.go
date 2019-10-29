package main

import (
	"math/big"
	"os"
	"strconv"
)

func parseKeySize(index *int) (uint, error) {
	(*index)++
	var argument string
	if *index < len(os.Args) {
		argument = os.Args[*index]
		if size, err := strconv.Atoi(argument); err == nil {
			return uint(size), nil
		}
	} else {
		argument = ""
	}
	return 0, &intFormatError{argument, keySizeBase, "key size"}
}

func parseKey(index *int) (*big.Int, error) {
	(*index)++
	var argument string
	var ok bool
	if *index < len(os.Args) {
		argument = os.Args[*index]
		key := new(big.Int)
		if key, ok = key.SetString(argument, keyBase); ok {
			return key, nil
		}
	} else {
		argument = ""
	}
	return nil, &intFormatError{argument, keyBase, "key"}
}
