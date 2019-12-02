package test

import (
	"bytes"
	"github.com/ikcilrep/twister/internal/functionality"
	"github.com/ikcilrep/twister/internal/parser"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func testEncryption(t *testing.T, data []byte, dataWriter *TestingWriter, key *big.Int, arguments *parser.Arguments) {
	err := functionality.Encrypt(data, key, arguments)
	if err != nil {
		t.Error(err)
	}
	encryptedData := dataWriter.data

	err = functionality.Decrypt(encryptedData, key, arguments)
	if err != nil {
		t.Error(err)
	}
	decryptedData := dataWriter.data

	if !bytes.Equal(data, decryptedData) {
		t.Errorf("%v is not %v", data, decryptedData)
	}
}

func Test_functionality_Encrypt(t *testing.T) {
	rand.Seed(time.Now().Unix())
	keySize := 32
	arguments := &parser.Arguments{KeySize: keySize, DataInput: &parser.Input{}}
	keyWriter := &TestingWriter{}
	dataWriter := &TestingWriter{}
	arguments.KeyOutput = &parser.Output{Writer: keyWriter, IsBinary: true}
	functionality.GenerateKey(arguments)
	arguments.KeyInput = &parser.Input{Reader: bytes.NewReader(keyWriter.data), IsBinary: true}
	arguments.DataOutput = &parser.Output{Writer: dataWriter}

	key, err := arguments.GetKey()
	if err != nil {
		t.Error(err)
	}

	for dataLength := 1; dataLength <= 128; dataLength++ {
		data := make([]byte, dataLength)
		rand.Read(data)
		testEncryption(t, data, dataWriter, key, arguments)
		arguments.DataOutput.IsBinary = true
		arguments.DataInput.IsBinary = true
		testEncryption(t, data, dataWriter, key, arguments)
	}
}
