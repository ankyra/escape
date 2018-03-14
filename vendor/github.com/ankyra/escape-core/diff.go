/*
Copyright 2017, 2018 Ankyra

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
	"reflect"
	"strconv"
	"strings"

	"github.com/ankyra/escape-core/templates"
	"github.com/ankyra/escape-core/variables"
)

type Change struct {
	Path          []string
	PreviousValue interface{}
	NewValue      interface{}
	Added         bool
	Removed       bool
}

func NewUpdate(path []string, old, new interface{}) Change {
	return Change{
		Path:          path,
		PreviousValue: old,
		NewValue:      new,
	}
}

func NewAddition(path []string, new interface{}) Change {
	return Change{
		Path:     path,
		NewValue: new,
		Added:    true,
	}
}

func NewRemoval(path []string, old interface{}) Change {
	return Change{
		Path:          path,
		PreviousValue: old,
		Removed:       true,
	}
}

func (c Change) ToString() string {
	name := strings.Join(c.Path, "")
	if !c.Added && !c.Removed {
		return fmt.Sprintf("Change %s from '%s' to '%s'", name, c.PreviousValue, c.NewValue)
	} else if c.Added {
		return fmt.Sprintf("Add '%s' to %s", c.NewValue, name)
	}
	return fmt.Sprintf("Remove '%s' from %s", c.PreviousValue, name)
}

func (c Change) GetModification() string {
	if !c.Added && !c.Removed {
		return "change"
	} else if c.Added {
		return "add"
	}
	return "remove"
}

type Changes []Change

func (changes Changes) Collapse() map[string]map[string]Changes {
	result := map[string]map[string]Changes{}
	for _, ch := range changes {
		field, exists := result[ch.Path[0]]
		if !exists {
			field = map[string]Changes{}
		}
		mod := ch.GetModification()
		modifications, exists := field[mod]
		if !exists {
			modifications = Changes{}
		}
		modifications = append(modifications, ch)
		field[mod] = modifications
		result[ch.Path[0]] = field
	}
	return result
}

func Diff(this *ReleaseMetadata, other *ReleaseMetadata) Changes {
	return diff([]string{}, this, other)
}

func diff(path []string, oldValue, newValue interface{}) Changes {
	if changes := diffNil(path, oldValue, newValue); len(changes) != 0 || oldValue == nil {
		return changes
	}
	thisVal := reflect.ValueOf(oldValue)
	typ := thisVal.Type().String()
	kind := thisVal.Type().Kind().String()
	if typ == "int" || typ == "string" || typ == "bool" {
		if r := diffSimpleType(path, oldValue, newValue); r != nil {
			return []Change{*r}
		}
	} else if kind == "ptr" {
		return diffPointer(path, oldValue, newValue)
	} else if kind == "struct" {
		return diffStruct(path, oldValue, newValue)
	} else if kind == "map" {
		return diffMap(path, oldValue, newValue)
	} else if kind == "slice" {
		return diffSlice(path, oldValue, newValue)
	} else {
		name := strings.Join(path, ", ")
		panic(fmt.Sprintf("WARN: Undiffable type '%s' (%s) for field '%s'\n", typ, kind, name))
	}
	return []Change{}
}

func diffStruct(path []string, oldValue, newValue interface{}) Changes {
	result := []Change{}
	oldVal := reflect.Indirect(reflect.ValueOf(oldValue))
	newVal := reflect.Indirect(reflect.ValueOf(newValue))
	fields := oldVal.Type().NumField()
	for i := 0; i < fields; i++ {
		field := oldVal.Type().Field(i).Name
		oldValue := oldVal.Field(i).Interface()
		newValue := newVal.FieldByName(field).Interface()
		newName := "." + field
		if len(path) == 0 {
			newName = field
		}
		structPath := make([]string, len(path)+1)
		copy(structPath, path)
		structPath[len(path)] = newName
		for _, change := range diff(structPath, oldValue, newValue) {
			result = append(result, change)
		}
	}
	return result
}

func diffSimpleType(path []string, oldValue, newValue interface{}) *Change {
	if !reflect.DeepEqual(oldValue, newValue) {
		v := NewUpdate(path, diffValue(oldValue), diffValue(newValue))
		return &v
	}
	return nil
}

func diffMap(path []string, oldValue, newValue interface{}) []Change {
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
			changes = append(changes, NewRemoval(path, key.String()))
			continue
		}
		newValI := newVal.Interface()
		if reflect.DeepEqual(oldVal, newValI) {
			continue
		}
		field := fmt.Sprintf(`["%s"]`, key)
		keyPath := make([]string, len(path)+1)
		copy(keyPath, path)
		keyPath[len(path)] = field
		for _, c := range diff(keyPath, oldVal, newValI) {
			changes = append(changes, c)
		}
	}
	for _, key := range newMap.MapKeys() {
		oldVal := oldMap.MapIndex(key)
		if !oldVal.IsValid() {
			changes = append(changes, NewAddition(path, key.String()))
		}
	}
	return changes
}
func diffSlice(path []string, oldValue, newValue interface{}) []Change {
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
				changes = append(changes, NewAddition(path, diffValue(val)))
			}
			continue
		} else if ix >= newValLen {
			if ix < oldValLen {
				val := oldVal.Index(ix).Interface()
				changes = append(changes, NewRemoval(path, diffValue(val)))
			}
			continue
		} else {
			ixPath := make([]string, len(path)+1)
			copy(ixPath, path)
			ixPath[len(path)] = "[" + strconv.Itoa(ix) + "]"
			for _, change := range diff(ixPath, oldVal.Index(ix).Interface(), newVal.Index(ix).Interface()) {
				changes = append(changes, change)
			}
		}
	}
	return changes
}

func diffNil(path []string, oldValue, newValue interface{}) []Change {
	changes := []Change{}
	if oldValue == nil {
		if newValue == nil {
			return changes
		} else {
			v := reflect.Indirect(reflect.ValueOf(newValue)).Interface()
			changes = append(changes, NewAddition(path, diffValue(v)))
		}
	} else if newValue == nil {
		v := reflect.Indirect(reflect.ValueOf(oldValue)).Interface()
		changes = append(changes, NewRemoval(path, diffValue(v)))
	}
	return changes
}

func diffPointer(path []string, oldValue, newValue interface{}) []Change {
	if changes := diffNil(path, oldValue, newValue); len(changes) != 0 {
		return changes
	}
	oldVal := reflect.Indirect(reflect.ValueOf(oldValue)).Interface()
	newVal := reflect.Indirect(reflect.ValueOf(newValue)).Interface()
	return diff(path, oldVal, newVal)
}

func diffValue(v interface{}) interface{} {
	switch v.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v.(int))
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
