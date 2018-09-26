package scopes

import "fmt"

type Scope string
type Scopes []string

const BuildScope = "build"
const DeployScope = "deploy"

var DeployScopes = Scopes{DeployScope}
var BuildScopes = Scopes{BuildScope}
var AllScopes = Scopes{BuildScope, DeployScope}

func NewScopesFromInterface(val interface{}) (Scopes, error) {
	result := Scopes{}
	valStr, ok := val.(string)
	if ok {
		result = append(result, valStr)
		return result.Validate()
	}
	valList, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Expecting string in scopes, got '%v' (%T)", val, val)
	}
	for _, val := range valList {
		kStr, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string in scopes, got '%v' (%T)", val, val)
		}
		result = append(result, kStr)
	}
	return result.Validate()
}

func (s Scopes) Validate() (Scopes, error) {
	for _, sc := range s {
		if sc != BuildScope && sc != DeployScope {
			return nil, fmt.Errorf("Unknown scope '%s'. Expecting '%s' or '%s'", sc, BuildScope, DeployScope)
		}
	}
	return s, nil
}

func (s Scopes) Copy() Scopes {
	result := make(Scopes, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = s[i]
	}
	return result
}

func (s Scopes) InScope(scope string) bool {
	for _, sc := range s {
		if sc == scope {
			return true
		}
	}
	return false
}
