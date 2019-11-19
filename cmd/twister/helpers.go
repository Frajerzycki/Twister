package main

import (
	"bufio"
	"io/ioutil"
	"os"
)

func getDataFromStd() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	return ioutil.ReadAll(reader)
}

func closeFiles(files []*os.File) {
	for _, v := range files {
		v.Close()
	}
}
