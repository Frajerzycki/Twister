package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"io/ioutil"
	"os"
)

var notPositiveLengthError error = errors.New("Length of data which exists has to be positive.")

func randomBytes(length int) ([]byte, error) {
	if length < 0 {
		return nil, notPositiveLengthError
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
