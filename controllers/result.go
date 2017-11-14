package controllers

import (
	"encoding/json"
	"fmt"
)

type ControllerResult struct {
	HumanOutput       *HumanOutput
	MarshalableOutput interface{}
	Error             error
}

func NewControllerResult() *ControllerResult {
	return &ControllerResult{
		HumanOutput: NewHumanOutput(""),
	}
}

func (r ControllerResult) Print(jsonFlag bool) error {
	if r.Error != nil {
		return r.Error
	}

	if jsonFlag {
		data, _ := json.Marshal(r.MarshalableOutput)
		fmt.Println(string(data))
	} else {
		fmt.Println(r.HumanOutput.value)
	}

	return nil
}

type HumanOutput struct {
	value string
}

func NewHumanOutput(format string, a ...interface{}) *HumanOutput {
	return &HumanOutput{
		value: fmt.Sprintf(format, a...),
	}
}

func (h *HumanOutput) AddLine(format string, a ...interface{}) {
	h.value = fmt.Sprintf(format, a...)
}

func (h *HumanOutput) AddMap(mapToAdd map[string]interface{}) {
	for k, v := range mapToAdd {
		h.value = NewHumanOutput("%s\n%s: %v", h.value, k, v).value
	}
}

func (h *HumanOutput) AddList(listToAdd []interface{}) {
	i := 0
	for _, k := range listToAdd {
		if i == 0 {
			h.value = NewHumanOutput("%v", k).value
		} else {
			h.value = NewHumanOutput("%s\n%v", h.value, k).value
		}
		i++
	}
}
