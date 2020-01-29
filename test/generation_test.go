package test

import (
	"bytes"
	"encoding/base64"
	"github.com/ikcilrep/twister/internal/functionality"
	"github.com/ikcilrep/twister/internal/parser"
	"testing"
)

func Test_functionality_GenerateKey(t *testing.T) {
	keyWriter := new(bytes.Buffer)
	arguments := &parser.Arguments{KeyWriter: keyWriter}
	for keySize := 0; keySize < 128; keySize++ {
		arguments.KeySize = keySize
		err := functionality.GenerateKey(arguments)
		if err != nil {
			t.Error(err)
		}
		decodedKeyBytes, err := base64.StdEncoding.DecodeString(string(keyWriter.Bytes()))
		if len(decodedKeyBytes) != keySize {
			t.Errorf("%v hasn't got length %v, but %v", keyWriter.Bytes(), keySize, len(keyWriter.Bytes()))
		}

		keyWriter.Reset()
	}
}
