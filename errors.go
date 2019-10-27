package main

import "fmt"

type keySizeFormatError struct {
	keySize string
}

func (err *keySizeFormatError) Error() string {
	return fmt.Sprintf("Wrong key size format, %v is not an integer.", err.keySize)
}
