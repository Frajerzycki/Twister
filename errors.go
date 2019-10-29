package main

import "fmt"

type intFormatError struct {
	representation string
	base           int
	name           string
}

func (err *intFormatError) Error() string {
	return fmt.Sprintf("Wrong %v format, \"%v\" is not an integer in %v base.", err.name, err.representation, err.base)
}
