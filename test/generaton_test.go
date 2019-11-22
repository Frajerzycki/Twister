package test

import (
	"encoding/base64"
	"github.com/ikcilrep/twister/internal/functionality"
	"github.com/ikcilrep/twister/internal/parser"
	"testing"
)

func TestKeyGenerating(t *testing.T) {
	keyWriter := &TestingWriter{}
	keyOutput := &parser.Output{Writer: keyWriter}
	arguments := &parser.Arguments{KeyOutput: keyOutput}
	for keySize := 0; keySize < 128; keySize++ {
		arguments.KeySize = keySize
		keyOutput.IsBinary = false
		err := functionality.GenerateKey(arguments)
		if err != nil {
			t.Error(err)
		}
		decodedKeyBytes, err := base64.StdEncoding.DecodeString(string(keyWriter.data))
		if len(decodedKeyBytes) != keySize {
			t.Errorf("%v hasn't got length %v, but %v", keyWriter.data, keySize, len(keyWriter.data))
		}

		keyOutput.IsBinary = true
		err = functionality.GenerateKey(arguments)
		if err != nil {
			t.Error(err)
		}
		if len(keyWriter.data) != keySize {
			t.Errorf("%v hasn't got length %v, but %v", keyWriter.data, keySize, len(keyWriter.data))
		}
	}
}
