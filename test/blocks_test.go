package test

import (
	"bytes"
	cryptoRand "crypto/rand"
	"github.com/ikcilrep/twister/internal/functionality"
	"math/rand"
	"testing"
)

func Test_functionality_ByteArrayToDataBlocks(t *testing.T) {
	for dataLength := 0; dataLength <= 128; dataLength++ {
		for blockLength := 1; blockLength <= dataLength; blockLength++ {
			data := make([]byte, dataLength)
			_, err := rand.Read(data)
			if err != nil {
				t.Error(err)
			}

			dataBlocks, err := functionality.ByteArrayToDataBlocks(data, blockLength, cryptoRand.Reader)
			if err != nil {
				t.Error(err)
			}

			retrievedData := dataBlocks.ToByteArray()
			if !bytes.Equal(data, retrievedData) {
				t.Errorf("%v\nis not the same as\n%v,\nrest is %v", retrievedData, data, dataBlocks.Rest)
			}
		}
	}
}
