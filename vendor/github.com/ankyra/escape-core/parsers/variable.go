package parsers

import (
	"fmt"
	"strings"
)

func InvalidVariableIdFormatError(v string) error {
	return fmt.Errorf("Invalid variable format '%s'", v)
}

func InvalidVariableIdPreviousError(v string) error {
	return fmt.Errorf("Invalid variable format '%s'. Variable is not allowed to start with 'PREVIOUS_'", v)
}

var InvalidVariableIdEmptyError = fmt.Errorf("Expecting variable string, but got empty string.")

func ParseVariableIdent(v string) (string, error) {
	if v == "" {
		return "", InvalidVariableIdEmptyError
	}
	variable, rest := ParseIdent(v)
	if strings.TrimSpace(rest) != "" {
		return "", InvalidVariableIdFormatError(v)
	}
	variable = strings.TrimSpace(variable)
	if strings.HasPrefix(strings.ToUpper(variable), "PREVIOUS_") {
		return "", InvalidVariableIdPreviousError(v)
	}
	return variable, nil
}
