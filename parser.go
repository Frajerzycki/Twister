package main

import (
	"math/big"
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
	if size, err := strconv.Atoi(argument); err == nil {
		return uint(size), nil
	}
	return 0, &intFormatError{argument, keySizeBase, "key size"}
}

func parseKey(index *int) (*big.Int, error) {
	argument := getCommandLineArgument(index)
	key := new(big.Int)
	var ok bool
	if key, ok = key.SetString(argument, keyBase); ok {
		return key, nil
	}
	return nil, &intFormatError{argument, keyBase, "key"}
}
