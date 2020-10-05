package parser

import (
	"fmt"
)

type manyParameterValuesError struct {
	parameterName string
}

func (err *manyParameterValuesError) Error() string {
	return fmt.Sprintf("%v has been set more than one time.", err.parameterName)
}
