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

func Test_files_ChunkedWriter(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for dataLength := 2; dataLength < 128; dataLength++ {
		data := make([]byte, dataLength)
		rand.Read(data)
		writer := new(bytes.Buffer)
		chunkedWriter, err := files.NewChunkedWriter(writer, 4)
		if err != nil {
			t.Error(err)
		}

		for _, v := range data {
			_, err = chunkedWriter.Write([]byte{v})
			if err != nil {
				t.Error(err)
			}

		}
		err = chunkedWriter.Close()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(data, writer.Bytes()) {
			t.Errorf("%v\nis not\n%v", data, writer.Bytes())
		}
	}
}
