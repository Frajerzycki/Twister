package test

import (
	"bytes"
	"github.com/ikcilrep/twister/internal/files"
	"io"
	"math/rand"
	"testing"
	"time"
)

func Test_files_ChunkedReader(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for dataLength := 2; dataLength < 128; dataLength++ {
		data := make([]byte, dataLength)
		rand.Read(data)
		reader := bytes.NewReader(data)
		chunkedReader, err := files.NewChunkedReader(reader, 4)
		if err != nil {
			t.Error(err)
		}
		readData := make([]byte, dataLength)
		for index := 0; err == nil && index+1 <= dataLength; index++ {
			_, err = io.ReadFull(chunkedReader, readData[index:index+1])
		}
		if !bytes.Equal(data, readData[:dataLength]) {
			t.Errorf("%v\nis not\n%v", data, readData)
		}
	}
}
