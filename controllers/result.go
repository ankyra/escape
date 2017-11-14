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

func (h *HumanOutput) AddLine(a string, a1 ...interface{}) {
	if h.addNewLine() {
		if len(a1) == 0 {
			h.value = fmt.Sprintf("%v\n%s", h.value, a)
		} else {
			a1 = append([]interface{}{h.value}, a1...)
			h.value = fmt.Sprintf("%v\n"+a, a1...)
		}

	} else {
		h.value = fmt.Sprintf(a, a1...)
	}
}

func (h *HumanOutput) AddMap(mapToAdd map[string]interface{}) {
	i := 0
	for k, v := range mapToAdd {
		if i == 0 {
			if h.addNewLine() {
				h.value = NewHumanOutput("%s\n\n%s: %v", h.value, k, v).value
			} else {
				h.value = NewHumanOutput("%s%s: %v", h.value, k, v).value
			}
		} else {
			h.value = NewHumanOutput("%s\n%s: %v", h.value, k, v).value
		}
		i++
	}
}

func (h *HumanOutput) AddList(listToAdd []interface{}) {
	i := 0
	for _, k := range listToAdd {
		if i == 0 {
			if h.addNewLine() {
				h.value = NewHumanOutput("%s\n\n%v", h.value, k).value
			} else {
				h.value = NewHumanOutput("%v", k).value
			}

		} else {
			h.value = NewHumanOutput("%s\n%v", h.value, k).value
		}
		i++
	}
}

func (h *HumanOutput) addNewLine() bool {
	return len(h.value) > 0
}
