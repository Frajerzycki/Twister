package main

import (
	"io/ioutil"
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
	size, err := strconv.Atoi(argument)
	if err != nil {
		return uint(size), nil
	}
	return 0, err
}

func parseKey(index *int) (*big.Int, error) {
	argument := getCommandLineArgument(index)
	key := new(big.Int)
	var ok bool
	if key, ok = key.SetString(argument, keyBase); ok {
		return key, nil
	}
	return nil, &formatError{argument, keyBase, "key", "integer"}
}

func parseKeyFromFile(index *int) (*big.Int, error) {
	argument := getCommandLineArgument(index)
	key := new(big.Int)
	file, err := os.Open(argument)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var keyAscii []byte
	if keyAscii, err = ioutil.ReadAll(file); err != nil {
		return nil, err
	}
	keyString := string(keyAscii)
	if key, ok := key.SetString(keyString, keyBase); ok {
		return key, nil
	}
	return nil, &formatError{keyString, keyBase, "key", "integer"}
}
