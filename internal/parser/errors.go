package parser

import (
	"errors"
	"fmt"
)

type manyParameterValuesError struct {
	parameterName string
}

func (err *manyParameterValuesError) Error() string {
	return fmt.Sprintf("%v has been set more than one time.", err.parameterName)
}

var binaryOutputError error = errors.New("Only output to file can be binary.")
var binaryInputError error = errors.New("Only input from file can be binary.")
