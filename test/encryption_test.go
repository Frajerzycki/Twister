package test

import (
	"bytes"
	"github.com/ikcilrep/twister/internal/functionality"
	"github.com/ikcilrep/twister/internal/parser"
	"math/rand"
	"testing"
	"time"
)

func Test_functionality_Encrypt(t *testing.T) {
	rand.Seed(time.Now().Unix())
	keySize := 32
	arguments := &parser.Arguments{KeySize: keySize}
	keyReaderWriter := NewTestBuffer()
	arguments.KeyWriter = keyReaderWriter
	arguments.KeyReader = keyReaderWriter
	functionality.GenerateKey(arguments)
	key, err := arguments.GetKey()
	if err != nil {
		t.Error(err)
	}

	for dataLength := 1; dataLength <= 128; dataLength++ {
		data := make([]byte, dataLength)
		rand.Read(data)
		arguments.DataReader = bytes.NewReader(data)
		buffer1 := NewTestBuffer()
		arguments.DataWriter = buffer1

		functionality.Encrypt(key, arguments)

		arguments.DataReader = buffer1
		buffer2 := NewTestBuffer()
		arguments.DataWriter = buffer2

		functionality.Decrypt(key, arguments)
		decryptedData := buffer2.Bytes()

		if !bytes.Equal(data, decryptedData) {
			t.Errorf("%v\nis not\n%v", data, decryptedData)
		}
	}
}
