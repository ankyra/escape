package controllers

import (
	"encoding/json"
	"fmt"
)

type ControllerResult struct {
	HumanOutput       string
	MarshalableOutput interface{}
	Error             error
}

func (r ControllerResult) Print(jsonFlag bool) error {
	if r.Error != nil {
		return r.Error
	}

	if jsonFlag {
		data, _ := json.Marshal(r.MarshalableOutput)
		fmt.Println(string(data))
	} else {
		fmt.Println(r.HumanOutput)
	}

	return nil
}
