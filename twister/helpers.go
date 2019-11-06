package main

import (
	"bufio"
	"crypto/rand"
	"io/ioutil"
	"os"
)

func randomBytes(length int) ([]byte, error) {
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
