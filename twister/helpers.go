package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"io/ioutil"
	"os"
)

var negativeLengthError error = errors.New("Length of any data musn't be negative.")

func randomBytes(length int) ([]byte, error) {
	if length < 0 {
		return nil, negativeLengthError
	}
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func getDataFromStd() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	return ioutil.ReadAll(reader)
}

func closeFiles(files []*os.File) {
	for _, v := range files {
		v.Close()
	}
}
