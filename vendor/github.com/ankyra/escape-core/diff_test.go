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
	"reflect"

	"github.com/ankyra/escape-core/templates"
	"github.com/ankyra/escape-core/variables"
	. "gopkg.in/check.v1"
)

func (s *metadataSuite) Test_Diff_simple_types(c *C) {
	testCases := [][]interface{}{
		[]interface{}{"ApiVersion", 1, 2, `Change ApiVersion from '1' to '2'`},
		[]interface{}{"Branch", "test", "not-test", `Change Branch from 'test' to 'not-test'`},
		[]interface{}{"Description", "test", "not-test", `Change Description from 'test' to 'not-test'`},
		[]interface{}{"Logo", "test", "not-test", `Change Logo from 'test' to 'not-test'`},
		[]interface{}{"Name", "test", "not-test", `Change Name from 'test' to 'not-test'`},
		[]interface{}{"Project", "test", "not-test", `Change Project from 'test' to 'not-test'`},
		[]interface{}{"Revision", "test", "not-test", `Change Revision from 'test' to 'not-test'`},
		[]interface{}{"Version", "1.0", "1.0.0", `Change Version from '1.0' to '1.0.0'`},
		[]interface{}{"Repository", "test", "not-test", `Change Repository from 'test' to 'not-test'`},
	}
	for _, test := range testCases {
		m1 := NewReleaseMetadata("test", "1.0")
		m2 := NewReleaseMetadata("test", "1.0")

		thisVal := reflect.Indirect(reflect.ValueOf(m1))
		otherVal := reflect.Indirect(reflect.ValueOf(m2))

		thisVal.FieldByName(test[0].(string)).Set(reflect.ValueOf(test[1]))
		otherVal.FieldByName(test[0].(string)).Set(reflect.ValueOf(test[2]))

		changes := Diff(m1, m2)
		c.Assert(changes, HasLen, 1, Commentf("Field %s", test[0]))
		c.Assert(changes[0].ToString(), Equals, test[3], Commentf("Field %s", test[0]))
	}
}

func (s *metadataSuite) Test_Diff_maps(c *C) {
	emptyDict := map[string]string{}
	oldDict := map[string]string{
		"newfile.txt": "123",
	}
	newDict := map[string]string{
		"newfile.txt": "123123123",
	}
	testCases := [][]interface{}{

		[]interface{}{"Files", oldDict, newDict, `Change Files["newfile.txt"] from '123' to '123123123'`},
		[]interface{}{"Files", emptyDict, newDict, `Add 'newfile.txt' to Files`},
		[]interface{}{"Files", oldDict, emptyDict, `Remove 'newfile.txt' from Files`},

		[]interface{}{"Metadata", oldDict, newDict, `Change Metadata["newfile.txt"] from '123' to '123123123'`},
		[]interface{}{"Metadata", emptyDict, newDict, `Add 'newfile.txt' to Metadata`},
		[]interface{}{"Metadata", oldDict, emptyDict, `Remove 'newfile.txt' from Metadata`},
	}
	for _, test := range testCases {
		m1 := NewReleaseMetadata("test", "1.0")
		m2 := NewReleaseMetadata("test", "1.0")

		thisVal := reflect.Indirect(reflect.ValueOf(m1))
		otherVal := reflect.Indirect(reflect.ValueOf(m2))

		thisVal.FieldByName(test[0].(string)).Set(reflect.ValueOf(test[1]))
		otherVal.FieldByName(test[0].(string)).Set(reflect.ValueOf(test[2]))

		changes := Diff(m1, m2)
		c.Assert(changes, HasLen, 1, Commentf(test[3].(string)))
		c.Assert(changes[0].ToString(), DeepEquals, test[3])
	}
}

func (s *metadataSuite) Test_Diff_Stages(c *C) {
	emptyDict := map[string]*ExecStage{}
	oldDict := map[string]*ExecStage{
		"test": &ExecStage{Script: "test.sh"},
	}
	newDict := map[string]*ExecStage{
		"test": &ExecStage{Script: "test2.sh"},
	}
	testCases := [][]interface{}{
		[]interface{}{oldDict, newDict, `Change Stages["test"].Script from 'test.sh' to 'test2.sh'`},
		[]interface{}{emptyDict, newDict, `Add 'test' to Stages`},
		[]interface{}{oldDict, emptyDict, `Remove 'test' from Stages`},
	}
	for _, test := range testCases {
		m1 := NewReleaseMetadata("test", "1.0")
		m2 := NewReleaseMetadata("test", "1.0")
		m1.Stages = test[0].(map[string]*ExecStage)
		m2.Stages = test[1].(map[string]*ExecStage)

		changes := Diff(m1, m2)
		c.Assert(changes, HasLen, 1, Commentf(test[2].(string)))
		c.Assert(changes[0].ToString(), DeepEquals, test[2])
	}
}

func (s *metadataSuite) Test_Diff_Errands(c *C) {
	errand1 := map[interface{}]interface{}{
		"script": "test.sh",
	}
	errand2 := map[interface{}]interface{}{
		"script": "test2.sh",
	}
	errand3 := map[interface{}]interface{}{
		"script":      "test.sh",
		"description": "Description",
	}

	testCases := [][]interface{}{
		[]interface{}{errand1, errand2, `Change Errands["test"].Script from 'test.sh' to 'test2.sh'`},
		[]interface{}{errand1, errand3, `Change Errands["test"].Description from '' to 'Description'`},
		[]interface{}{nil, errand1, `Add 'test' to Errands`},
		[]interface{}{errand1, nil, `Remove 'test' from Errands`},
	}
	for _, test := range testCases {
		m1 := NewReleaseMetadata("test", "1.0")
		m2 := NewReleaseMetadata("test", "1.0")

		if test[0] != nil {
			e1, err := NewErrandFromDict("test", test[0].(map[interface{}]interface{}))
			c.Assert(err, IsNil)
			m1.Errands["test"] = e1
		}
		if test[1] != nil {
			e2, err2 := NewErrandFromDict("test", test[1].(map[interface{}]interface{}))
			c.Assert(err2, IsNil)
			m2.Errands["test"] = e2
		}
		changes := Diff(m1, m2)
		c.Assert(changes, HasLen, 1, Commentf(test[2].(string)))
		c.Assert(changes[0].ToString(), DeepEquals, test[2])
	}
}

func (s *metadataSuite) Test_Diff_Variables(c *C) {
	empty := []map[interface{}]interface{}{}
	var1 := []map[interface{}]interface{}{
		map[interface{}]interface{}{
			"id": "test",
		},
	}
	var2 := []map[interface{}]interface{}{
		map[interface{}]interface{}{
			"id": "test2",
		},
	}
	var3 := []map[interface{}]interface{}{
		map[interface{}]interface{}{
			"id":   "test",
			"type": "integer",
		},
	}
	var4 := []map[interface{}]interface{}{
		map[interface{}]interface{}{
			"id": "test2",
		},
		map[interface{}]interface{}{
			"id": "test",
		},
	}

	testCases := [][]interface{}{
		[]interface{}{"Inputs", var1, var2, `Change Inputs[0].Id from 'test' to 'test2'`},
		[]interface{}{"Inputs", var1, var3, `Change Inputs[0].Type from 'string' to 'integer'`},
		[]interface{}{"Inputs", var1, empty, `Remove 'test' from Inputs`},
		[]interface{}{"Inputs", empty, var1, `Add 'test' to Inputs`},
		[]interface{}{"Inputs", var4, var2, `Remove 'test' from Inputs`},
		[]interface{}{"Inputs", var2, var4, `Add 'test' to Inputs`},

		[]interface{}{"Outputs", var1, var2, `Change Outputs[0].Id from 'test' to 'test2'`},
		[]interface{}{"Outputs", var1, var3, `Change Outputs[0].Type from 'string' to 'integer'`},
		[]interface{}{"Outputs", var1, empty, `Remove 'test' from Outputs`},
		[]interface{}{"Outputs", empty, var1, `Add 'test' to Outputs`},
		[]interface{}{"Outputs", var4, var2, `Remove 'test' from Outputs`},
		[]interface{}{"Outputs", var2, var4, `Add 'test' to Outputs`},

		[]interface{}{"Errands", var1, var2, `Change Errands["test"].Inputs[0].Id from 'test' to 'test2'`},
		[]interface{}{"Errands", var1, var3, `Change Errands["test"].Inputs[0].Type from 'string' to 'integer'`},
		[]interface{}{"Errands", var1, empty, `Remove 'test' from Errands["test"].Inputs`},
		[]interface{}{"Errands", empty, var1, `Add 'test' to Errands["test"].Inputs`},
		[]interface{}{"Errands", var4, var2, `Remove 'test' from Errands["test"].Inputs`},
		[]interface{}{"Errands", var2, var4, `Add 'test' to Errands["test"].Inputs`},
	}
	for _, test := range testCases {
		errand1, err := NewErrandFromDict("test", map[interface{}]interface{}{
			"script": "test.sh",
		})
		c.Assert(err, IsNil)
		errand2, err := NewErrandFromDict("test", map[interface{}]interface{}{
			"script": "test.sh",
		})
		c.Assert(err, IsNil)
		m1 := NewReleaseMetadata("test", "1.0")
		m2 := NewReleaseMetadata("test", "1.0")
		m1.Errands["test"] = errand1
		m2.Errands["test"] = errand2
		typ := test[0].(string)
		for _, varDict := range test[1].([]map[interface{}]interface{}) {
			v, err := variables.NewVariableFromDict(varDict)
			c.Assert(err, IsNil)
			if typ == "Inputs" {
				m1.AddInputVariable(v)
			} else if typ == "Outputs" {
				m1.AddOutputVariable(v)
			} else {
				errand1.Inputs = append(errand1.Inputs, v)
			}
		}
		for _, varDict := range test[2].([]map[interface{}]interface{}) {
			v, err := variables.NewVariableFromDict(varDict)
			c.Assert(err, IsNil)
			if typ == "Inputs" {
				m2.AddInputVariable(v)
			} else if typ == "Outputs" {
				m2.AddOutputVariable(v)
			} else {
				errand2.Inputs = append(errand2.Inputs, v)
			}
		}

		changes := Diff(m1, m2)
		c.Assert(changes, HasLen, 1, Commentf(test[3].(string)))
		c.Assert(changes[0].ToString(), DeepEquals, test[3])
	}
}

func (s *metadataSuite) Test_Diff_Slices(c *C) {
	testCases := [][]interface{}{
		[]interface{}{"Provides", []string{"test"}, []string{}, `Remove 'test' from Provides`},
		[]interface{}{"Provides", []string{}, []string{"test"}, `Add 'test' to Provides`},
		[]interface{}{"Provides", []string{"test"}, []string{"kubernetes"}, `Change Provides[0].Name from 'test' to 'kubernetes'`},
		[]interface{}{"Extends", []string{"test"}, []string{}, `Remove 'test' from Extends`},
		[]interface{}{"Extends", []string{}, []string{"test"}, `Add 'test' to Extends`},
		[]interface{}{"Extends", []string{"test"}, []string{"kubernetes"}, `Change Extends[0].ReleaseId from 'test' to 'kubernetes'`},
	}
	for _, test := range testCases {
		m1 := NewReleaseMetadata("test", "1.0")
		m2 := NewReleaseMetadata("test", "1.0")
		typ := test[0].(string)
		if typ == "Consumes" {
			m1.SetConsumes(test[1].([]string))
			m2.SetConsumes(test[2].([]string))
		} else if typ == "Provides" {
			m1.SetProvides(test[1].([]string))
			m2.SetProvides(test[2].([]string))
		} else if typ == "Depends" {
			m1.SetDependencies(test[1].([]string))
			m2.SetDependencies(test[2].([]string))
		} else if typ == "Extends" {
			for _, consumer := range test[1].([]string) {
				m1.AddExtension(consumer)
			}
			for _, consumer := range test[2].([]string) {
				m2.AddExtension(consumer)
			}
		}
		changes := Diff(m1, m2)
		c.Assert(changes, HasLen, 1, Commentf(test[3].(string)))
		c.Assert(changes[0].ToString(), DeepEquals, test[3])
	}
}

func (s *metadataSuite) Test_Diff_Depends(c *C) {
	testCases := [][]interface{}{
		[]interface{}{"Depends", []string{"test-v1.0"}, []string{}, `Remove '_/test-v1.0' from Depends`, 1},
		[]interface{}{"Depends", []string{}, []string{"test-v1.0"}, `Add '_/test-v1.0' to Depends`, 1},
		[]interface{}{"Depends", []string{"test-v1.0"}, []string{"kubernetes-v1.0"}, `Change Depends[0].ReleaseId from '_/test-v1.0' to '_/kubernetes-v1.0'`, 4},
	}
	for _, test := range testCases {
		m1 := NewReleaseMetadata("test", "1.0")
		m2 := NewReleaseMetadata("test", "1.0")
		m1.SetDependencies(test[1].([]string))
		m2.SetDependencies(test[2].([]string))
		changes := Diff(m1, m2)
		c.Assert(changes, HasLen, test[4], Commentf(test[3].(string)))
		c.Assert(changes[0].ToString(), DeepEquals, test[3])
	}
}

func (s *metadataSuite) Test_Diff_Consumes(c *C) {
	testCases := [][]interface{}{
		[]interface{}{"Consumes", []string{"test"}, []string{}, `Remove 'test' from Consumes`, 1},
		[]interface{}{"Consumes", []string{}, []string{"test"}, `Add 'test' to Consumes`, 1},
		[]interface{}{"Consumes", []string{"test"}, []string{"kubernetes"}, `Change Consumes[0].Name from 'test' to 'kubernetes'`, 2, `Change Consumes[0].VariableName from 'test' to 'kubernetes'`},
	}
	for _, test := range testCases {
		m1 := NewReleaseMetadata("test", "1.0")
		m2 := NewReleaseMetadata("test", "1.0")
		m1.SetConsumes(test[1].([]string))
		m2.SetConsumes(test[2].([]string))
		changes := Diff(m1, m2)
		c.Assert(changes, HasLen, test[4].(int), Commentf(test[3].(string)))
		c.Assert(changes[0].ToString(), DeepEquals, test[3])
		if test[4].(int) > 1 {
			c.Assert(changes[1].ToString(), DeepEquals, test[5])
		}
	}
}

func (s *metadataSuite) Test_Diff_Templates(c *C) {
	v1 := map[interface{}]interface{}{
		"file":   "test.tpl",
		"target": "test.sh",
	}
	v2 := map[interface{}]interface{}{
		"file":   "test2.tpl",
		"target": "test.sh",
	}
	v3 := map[interface{}]interface{}{
		"file":   "test.tpl",
		"target": "/other",
	}
	testCases := [][]interface{}{
		[]interface{}{v1, v2, `Change Templates[0].File from 'test.tpl' to 'test2.tpl'`},
		[]interface{}{v1, v3, `Change Templates[0].Target from 'test.sh' to '/other'`},
		[]interface{}{nil, v1, `Add 'test.tpl' to Templates`},
		[]interface{}{v1, nil, `Remove 'test.tpl' from Templates`},
	}
	for _, test := range testCases {
		m1 := NewReleaseMetadata("test", "1.0")
		m2 := NewReleaseMetadata("test", "1.0")

		if test[0] != nil {
			e1, err := templates.NewTemplateFromInterface(test[0].(map[interface{}]interface{}))
			c.Assert(err, IsNil)
			m1.Templates = append(m1.Templates, e1)
		}

		if test[1] != nil {
			e2, err2 := templates.NewTemplateFromInterface(test[1].(map[interface{}]interface{}))
			c.Assert(err2, IsNil)
			m2.Templates = append(m2.Templates, e2)
		}
		changes := Diff(m1, m2)
		c.Assert(changes, HasLen, 1, Commentf(test[2].(string)))
		c.Assert(changes[0].ToString(), DeepEquals, test[2])
	}
}
