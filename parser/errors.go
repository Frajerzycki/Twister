package parser

import (
	"errors"
	"fmt"
)

type formatError struct {
	representation string
	base           int
	name           string
	dataType       string
}

type manyParameterValuesError struct {
	parameterName string
}

func (err *formatError) Error() string {
	return fmt.Sprintf("Wrong %v format, \"%v\" is not an %v in %v base.", err.name, err.representation, err.dataType, err.base)
}

func (err *manyParameterValuesError) Error() string {
	return fmt.Sprintf("%v has been set more than one time.", err.parameterName)
}

var binaryOutputError error = errors.New("Only output to file can be binary.")
var binaryInputError error = errors.New("Only input from file can be binary.")
