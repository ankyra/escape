/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"fmt"
	"github.com/ankyra/escape-core/templates"
	"github.com/ankyra/escape-core/variables"
	"reflect"
	"strconv"
)

type Change struct {
	Field         string
	PreviousValue interface{}
	NewValue      interface{}
	Added         bool
	Removed       bool
}

func NewUpdate(field string, old, new interface{}) Change {
	return Change{
		Field:         field,
		PreviousValue: old,
		NewValue:      new,
	}
}

func NewAddition(field string, new interface{}) Change {
	return Change{
		Field:    field,
		NewValue: new,
		Added:    true,
	}
}

func NewRemoval(field string, old interface{}) Change {
	return Change{
		Field:         field,
		PreviousValue: old,
		Removed:       true,
	}
}

func (c Change) ToString() string {
	if !c.Added && !c.Removed {
		return fmt.Sprintf("Change %s from '%s' to '%s'", c.Field, c.PreviousValue, c.NewValue)
	} else if c.Added {
		return fmt.Sprintf("Add '%s' to %s", c.NewValue, c.Field)
	}
	return fmt.Sprintf("Remove '%s' from %s", c.PreviousValue, c.Field)
}

type Changes []Change

func Diff(this *ReleaseMetadata, other *ReleaseMetadata) Changes {
	return diff("", this, other)
}

func diff(name string, oldValue, newValue interface{}) Changes {
	if changes := diffNil(name, oldValue, newValue); len(changes) != 0 || oldValue == nil {
		return changes
	}
	thisVal := reflect.ValueOf(oldValue)
	typ := thisVal.Type().String()
	kind := thisVal.Type().Kind().String()
	if typ == "int" || typ == "string" || typ == "bool" {
		if r := diffSimpleType(name, oldValue, newValue); r != nil {
			return []Change{*r}
		}
	} else if kind == "ptr" {
		return diffPointer(name, oldValue, newValue)
	} else if kind == "struct" {
		return diffStruct(name, oldValue, newValue)
	} else if kind == "map" {
		return diffMap(name, oldValue, newValue)
	} else if kind == "slice" {
		return diffSlice(name, oldValue, newValue)
	} else {
		panic(fmt.Sprintf("WARN: Undiffable type '%s' (%s) for field '%s'\n", typ, kind, name))
	}
	return []Change{}
}

func diffStruct(name string, oldValue, newValue interface{}) Changes {
	result := []Change{}
	oldVal := reflect.Indirect(reflect.ValueOf(oldValue))
	newVal := reflect.Indirect(reflect.ValueOf(newValue))
	fields := oldVal.Type().NumField()
	for i := 0; i < fields; i++ {
		field := oldVal.Type().Field(i).Name
		oldValue := oldVal.Field(i).Interface()
		newValue := newVal.FieldByName(field).Interface()
		newName := name + "." + field
		if name == "" {
			newName = field
		}
		if oldVal.NumField() == 1 {
			newName = name
		}
		for _, change := range diff(newName, oldValue, newValue) {
			result = append(result, change)
		}
	}
	return result
}

func diffSimpleType(name string, oldValue, newValue interface{}) *Change {
	if !reflect.DeepEqual(oldValue, newValue) {
		v := NewUpdate(name, diffValue(oldValue), diffValue(newValue))
		return &v
	}
	return nil
}

func diffMap(name string, oldValue, newValue interface{}) []Change {
	changes := []Change{}
	if reflect.DeepEqual(oldValue, newValue) {
		return changes
	}
	oldMap := reflect.ValueOf(oldValue)
	newMap := reflect.ValueOf(newValue)

	for _, key := range oldMap.MapKeys() {
		oldVal := oldMap.MapIndex(key).Interface()
		newVal := newMap.MapIndex(key)
		if !newVal.IsValid() {
			changes = append(changes, NewRemoval(name, key))
			continue
		}
		newValI := newVal.Interface()
		if reflect.DeepEqual(oldVal, newValI) {
			continue
		}
		field := fmt.Sprintf(`%s["%s"]`, name, key)
		for _, c := range diff(field, oldVal, newValI) {
			changes = append(changes, c)
		}
	}
	for _, key := range newMap.MapKeys() {
		oldVal := oldMap.MapIndex(key)
		if !oldVal.IsValid() {
			changes = append(changes, NewAddition(name, key))
		}
	}
	return changes
}
func diffSlice(name string, oldValue, newValue interface{}) []Change {
	changes := []Change{}
	if reflect.DeepEqual(oldValue, newValue) {
		return nil
	}
	oldVal := reflect.ValueOf(oldValue)
	oldValLen := oldVal.Len()
	newVal := reflect.ValueOf(newValue)
	newValLen := newVal.Len()
	until := oldValLen
	if newValLen > oldValLen {
		until = newValLen
	}
	for ix := 0; ix < until; ix++ {
		if ix >= oldValLen {
			if ix < newValLen {
				val := newVal.Index(ix).Interface()
				changes = append(changes, NewAddition(name, diffValue(val)))
			}
			continue
		} else if ix >= newValLen {
			if ix < oldValLen {
				val := oldVal.Index(ix).Interface()
				changes = append(changes, NewRemoval(name, diffValue(val)))
			}
			continue
		} else {
			for _, change := range diff(name+"["+strconv.Itoa(ix)+"]", oldVal.Index(ix).Interface(), newVal.Index(ix).Interface()) {
				changes = append(changes, change)
			}
		}
	}
	return changes
}

func diffNil(name string, oldValue, newValue interface{}) []Change {
	changes := []Change{}
	if oldValue == nil {
		if newValue == nil {
			return changes
		} else {
			v := reflect.Indirect(reflect.ValueOf(newValue)).Interface()
			changes = append(changes, Change{name, nil, diffValue(v), true, false})
		}
	} else if newValue == nil {
		v := reflect.Indirect(reflect.ValueOf(oldValue)).Interface()
		changes = append(changes, Change{name, diffValue(v), nil, true, false})
	}
	return changes
}

func diffPointer(name string, oldValue, newValue interface{}) []Change {
	if changes := diffNil(name, oldValue, newValue); len(changes) != 0 {
		return changes
	}
	oldVal := reflect.Indirect(reflect.ValueOf(oldValue)).Interface()
	newVal := reflect.Indirect(reflect.ValueOf(newValue)).Interface()
	return diff(name, oldVal, newVal)
}

func diffValue(v interface{}) interface{} {
	switch v.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v.(int))
	case *ExecStage:
		return v.(*ExecStage).Script
	case *ConsumerConfig:
		return v.(*ConsumerConfig).Name
	case *ProviderConfig:
		return v.(*ProviderConfig).Name
	case *DependencyConfig:
		return v.(*DependencyConfig).ReleaseId
	case *ExtensionConfig:
		return v.(*ExtensionConfig).ReleaseId
	case *variables.Variable:
		return v.(*variables.Variable).Id
	case *templates.Template:
		return v.(*templates.Template).File
	}
	return v
}
