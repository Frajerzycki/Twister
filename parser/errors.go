package parser

import "fmt"

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
