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

package escape_plan

import (
	. "gopkg.in/check.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strings"
)

type printSuite struct{}

var _ = Suite(&printSuite{})

func (s *printSuite) Test_PrettyPrint_KeyVal(c *C) {
	unit := NewPrettyPrinter(includeDocs(false), includeEmpty(true))
	keys := []string{"name", "version", "description", "logo", "path",
		"pre_build", "build", "post_build", "pre_destroy", "destroy", "post_destroy", "test",
		"pre_deploy", "deploy", "post_deploy", "smoke"}
	for _, key := range keys {
		pretty := unit.prettyPrintValue(key, "value")
		c.Assert(string(pretty), Equals, key+": value")
		pretty = unit.prettyPrintValue(key, "")
		c.Assert(string(pretty), Equals, key+": \"\"")
		pretty = unit.prettyPrintValue(key, nil)
		c.Assert(string(pretty), Equals, key+": \"\"")
	}
}

func (s *printSuite) Test_PrettyPrint_ListVal(c *C) {
	unit := NewPrettyPrinter(includeDocs(false), includeEmpty(true))
	keys := []string{"depends", "consumes", "provides", "includes", "inputs", "outputs", "templates"}
	for _, key := range keys {
		pretty := unit.prettyPrintValue(key, []string{})
		c.Assert(string(pretty), Equals, key+": []")
		pretty = unit.prettyPrintValue(key, nil)
		c.Assert(string(pretty), Equals, key+": []")
		pretty = unit.prettyPrintValue(key, []string{"test", "val"})
		c.Assert(string(pretty), Equals, key+":\n- test\n- val")
	}
}

func (s *printSuite) Test_PrettyPrint_MapVal(c *C) {
	unit := NewPrettyPrinter(includeDocs(false), includeEmpty(true))
	keys := []string{"metadata", "errands"}
	for _, key := range keys {
		pretty := unit.prettyPrintValue(key, map[string]string{})
		c.Assert(string(pretty), Equals, key+": {}")
		testMap := map[string]string{"test": "val", "test2": "val2"}
		pretty = unit.prettyPrintValue(key, testMap)
		c.Assert(string(pretty), Equals, key+":\n  test: val\n  test2: val2")
	}
}

func (s *printSuite) Test_PrettyPrint_Full_Fixture_No_Doc(c *C) {
	unit := NewPrettyPrinter(includeDocs(false), spacing(1))
	c.Assert(unit.IncludeEmpty, Equals, true)
	c.Assert(unit.IncludeDocs, Equals, false)
	c.Assert(unit.Spacing, Equals, 1)
	plan := NewEscapePlan()
	err := plan.LoadConfig("testdata/fixture.yml")
	c.Assert(err, IsNil)
	pretty := unit.Print(plan)
	expected, err := ioutil.ReadFile("testdata/fixture_pretty_print_no_doc.yml")
	c.Assert(err, IsNil)
	c.Assert(string(pretty), Equals, string(expected))
}

func (s *printSuite) Test_PrettyPrint_Minify_Full_Fixture_Same_As_No_Doc(c *C) {
	unit := NewPrettyPrinter(includeEmpty(false), includeDocs(false), spacing(1))
	c.Assert(unit.IncludeEmpty, Equals, false)
	c.Assert(unit.IncludeDocs, Equals, false)
	c.Assert(unit.Spacing, Equals, 1)
	plan := NewEscapePlan()
	err := plan.LoadConfig("testdata/fixture.yml")
	c.Assert(err, IsNil)
	pretty := unit.Print(plan)
	expected, err := ioutil.ReadFile("testdata/fixture_pretty_print_no_doc.yml")
	c.Assert(err, IsNil)
	c.Assert(string(pretty), Equals, string(expected))
}

func (s *printSuite) Test_PrettyPrint_Full_Fixture(c *C) {
	unit := NewPrettyPrinter()
	c.Assert(unit.IncludeEmpty, Equals, true)
	c.Assert(unit.IncludeDocs, Equals, true)
	c.Assert(unit.Spacing, Equals, 2)
	plan := NewEscapePlan()
	err := plan.LoadConfig("testdata/fixture.yml")
	c.Assert(err, IsNil)
	pretty := unit.Print(plan)
	expected, err := ioutil.ReadFile("testdata/fixture_pretty_print.yml")
	c.Assert(err, IsNil)
	c.Assert(string(pretty), Equals, string(expected))
}

func (s *printSuite) Test_PrettyPrint_Minimal_Fixture(c *C) {
	unit := NewPrettyPrinter(includeEmpty(false))
	c.Assert(unit.IncludeEmpty, Equals, false)
	c.Assert(unit.IncludeDocs, Equals, true)
	c.Assert(unit.Spacing, Equals, 2)
	plan := NewEscapePlan()
	err := plan.LoadConfig("testdata/minimal_fixture.yml")
	c.Assert(err, IsNil)
	pretty := unit.Print(plan)
	expected, err := ioutil.ReadFile("testdata/minimal_fixture_minified.yml")
	c.Assert(err, IsNil)
	c.Assert(string(pretty), Equals, string(expected))
}

func (s *printSuite) Test_PrettyPrint_includes_all_fields(c *C) {
	unit := NewPrettyPrinter(includeEmpty(true), includeDocs(true))
	plan := NewEscapePlan()
	pretty := unit.Print(plan)
	result := map[string]interface{}{}
	err := yaml.Unmarshal(pretty, &result)
	c.Assert(err, IsNil)
	val := reflect.Indirect(reflect.ValueOf(plan))
	for i := 0; i < val.Type().NumField(); i++ {
		name := val.Type().Field(i).Tag.Get("yaml")
		key := strings.Split(name, ",")[0]
		if key != "-" {
			_, found := result[key]
			c.Assert(found, Equals, true)
		}
	}
}
