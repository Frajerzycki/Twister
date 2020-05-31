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
	err := functionality.GenerateKey(arguments)
	if err != nil {
		t.Error(err)
	}

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

		_, _, err = functionality.Encrypt(key, arguments)
		if err != nil {
			t.Error(err)
		}

		arguments.DataReader = buffer1
		buffer2 := NewTestBuffer()
		arguments.DataWriter = buffer2

		_, _, err = functionality.Decrypt(key, arguments)
		if err != nil {
			t.Error(err)
		}
		decryptedData := buffer2.Bytes()

		if !bytes.Equal(data, decryptedData) {
			t.Errorf("%v\nis not\n%v", data, decryptedData)
		}
	}
}
