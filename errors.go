package main

import "fmt"

type intFormatError struct {
	representation string
	base           int
}

func (err *intFormatError) Error() string {
	return fmt.Sprintf("Wrong key size format, %v is not an integer in %v base.", err.representation, err.base)
}
