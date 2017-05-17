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

package templates

import (
	"github.com/ankyra/escape-core/script"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"testing"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *testSuite) Test_Template_renderToString(c *C) {
	mapping := map[string]interface{}{
		"who": "world",
	}
	unit := NewTemplateWithMapping("testdata/helloworld.mustache", mapping)
	env := script.NewScriptEnvironmentWithGlobals(nil)
	str, err := unit.renderToString(env)
	c.Assert(err, IsNil)
	c.Assert(str, Equals, "Hello world\n")
}

func (s *testSuite) Test_Template_renderToString_fails_if_mapping_is_missing(c *C) {
	mapping := map[string]interface{}{}
	unit := NewTemplateWithMapping("testdata/helloworld.mustache", mapping)
	env := script.NewScriptEnvironmentWithGlobals(nil)
	_, err := unit.renderToString(env)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_Template_renderToString_with_int(c *C) {
	mapping := map[string]interface{}{
		"who": 12,
	}
	unit := NewTemplateWithMapping("testdata/helloworld.mustache", mapping)
	env := script.NewScriptEnvironmentWithGlobals(nil)
	str, err := unit.renderToString(env)
	c.Assert(err, IsNil)
	c.Assert(str, Equals, "Hello 12\n")
}

func (s *testSuite) Test_Template_RenderScriptToString(c *C) {
	mapping := map[string]interface{}{
		"who": "$this.variable",
	}
	unit := NewTemplateWithMapping("testdata/helloworld.mustache", mapping)
	thisDict := map[string]script.Script{
		"variable": script.LiftString("scripted world"),
	}
	globalsDict := map[string]script.Script{
		"this": script.LiftDict(thisDict),
	}
	env := script.NewScriptEnvironmentWithGlobals(globalsDict)
	str, err := unit.renderToString(env)
	c.Assert(err, IsNil)
	c.Assert(str, Equals, "Hello scripted world\n")
}

func (s *testSuite) Test_Template_RenderScriptToString_fails_if_parse_fails(c *C) {
	mapping := map[string]interface{}{
		"who": "$this$doesnt$parse",
	}
	unit := NewTemplateWithMapping("testdata/helloworld.mustache", mapping)
	env := script.NewScriptEnvironmentWithGlobals(nil)
	_, err := unit.renderToString(env)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_Template_RenderScriptToString_fails_if_eval_fails(c *C) {
	mapping := map[string]interface{}{
		"who": "$this.doesnt_exist",
	}
	unit := NewTemplateWithMapping("testdata/helloworld.mustache", mapping)
	env := script.NewScriptEnvironmentWithGlobals(nil)
	_, err := unit.renderToString(env)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_file(c *C) {
	dict := map[string]interface{}{
		"file": "test.sh.tpl",
	}
	unit, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, IsNil)
	c.Assert(unit.File, Equals, "test.sh.tpl")
	c.Assert(unit.Target, Equals, "test.sh")
	c.Assert(unit.Mapping, HasLen, 0)
	c.Assert(unit.Scopes, DeepEquals, []string{"build", "deploy"})
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_file_fails_on_wrong_type(c *C) {
	dict := map[string]interface{}{
		"file": 20,
	}
	_, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_target(c *C) {
	dict := map[string]interface{}{
		"target": "test.sh",
	}
	unit, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, IsNil)
	c.Assert(unit.File, Equals, "")
	c.Assert(unit.Target, Equals, "test.sh")
	c.Assert(unit.Mapping, HasLen, 0)
	c.Assert(unit.Scopes, DeepEquals, []string{"build", "deploy"})
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_target_fails_on_wrong_type(c *C) {
	dict := map[string]interface{}{
		"target": 20,
	}
	_, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_mapping(c *C) {
	dict := map[string]interface{}{
		"mapping": map[interface{}]interface{}{
			"variable": "$this.inputs.hello",
		},
	}
	unit, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, IsNil)
	c.Assert(unit.File, Equals, "")
	c.Assert(unit.Target, Equals, "")
	c.Assert(unit.Mapping, HasLen, 1)
	c.Assert(unit.Mapping["variable"], Equals, "$this.inputs.hello")
	c.Assert(unit.Scopes, DeepEquals, []string{"build", "deploy"})
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_mapping_fails_on_wrong_type(c *C) {
	dict := map[string]interface{}{
		"mapping": 20,
	}
	_, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_mapping_fails_on_wrong_key_type(c *C) {
	dict := map[string]interface{}{
		"mapping": map[interface{}]interface{}{
			20: "$this.inputs.hello",
		},
	}
	_, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_scopes(c *C) {
	dict := map[string]interface{}{
		"scopes": []interface{}{"build", "some_stage"},
	}
	unit, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, IsNil)
	c.Assert(unit.File, Equals, "")
	c.Assert(unit.Target, Equals, "")
	c.Assert(unit.Mapping, HasLen, 0)
	c.Assert(unit.Scopes, DeepEquals, []string{"build", "some_stage"})
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_scope_as_string(c *C) {
	dict := map[string]interface{}{
		"scopes": "build",
	}
	unit, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, IsNil)
	c.Assert(unit.File, Equals, "")
	c.Assert(unit.Target, Equals, "")
	c.Assert(unit.Mapping, HasLen, 0)
	c.Assert(unit.Scopes, DeepEquals, []string{"build"})
}

func (s *testSuite) Test_NewTemplateFromInterfaceMap_scopes_fails_on_wrong_type(c *C) {
	dict := map[string]interface{}{
		"scopes": 20,
	}
	_, err := NewTemplateFromInterfaceMap(dict)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_NewTemplateFromInterface_from_string(c *C) {
	unit, err := NewTemplateFromInterface("testfile.txt")
	c.Assert(err, IsNil)
	c.Assert(unit.File, Equals, "testfile.txt")
	c.Assert(unit.Target, Equals, "testfile")
	c.Assert(unit.Mapping, HasLen, 0)
	c.Assert(unit.Scopes, DeepEquals, []string{"build", "deploy"})
}

func (s *testSuite) Test_NewTemplateFromInterface_from_map_fails_if_key_not_string(c *C) {
	dict := map[interface{}]interface{}{
		20: "test.sh.tpl",
	}
	_, err := NewTemplateFromInterface(dict)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_NewTemplateFromInterface_from_map(c *C) {
	dict := map[interface{}]interface{}{
		"file":   "test.sh.tpl",
		"target": "result.sh",
		"mapping": map[interface{}]interface{}{
			"variable": "$this.inputs",
		},
		"scopes": []interface{}{"lol", "test"},
	}
	unit, err := NewTemplateFromInterface(dict)
	c.Assert(err, IsNil)
	c.Assert(unit.File, Equals, "test.sh.tpl")
	c.Assert(unit.Target, Equals, "result.sh")
	c.Assert(unit.Mapping, HasLen, 1)
	c.Assert(unit.Mapping["variable"], Equals, "$this.inputs")
	c.Assert(unit.Scopes, DeepEquals, []string{"lol", "test"})
}

func (s *testSuite) Test_NewTemplateFromInterface_fails_with_unknown_type(c *C) {
	_, err := NewTemplateFromInterface(12)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_Render_fails_if_file_empty(c *C) {
	unit := NewTemplate()
	env := script.NewScriptEnvironmentWithGlobals(nil)
	err := unit.Render("build", env)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Can't run template. Template file has not been defined (missing 'file' key in Escape plan?)")
}

func (s *testSuite) Test_Render_fails_if_target_empty(c *C) {
	unit := NewTemplate()
	unit.SetFile("file.template")
	env := script.NewScriptEnvironmentWithGlobals(nil)
	err := unit.Render("build", env)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "Can't run template. Template target has not been defined (empty 'target' key in Escape plan?)")
}

func (s *testSuite) Test_Render_doesnt_run_if_not_in_scope(c *C) {
	unit := NewTemplate()
	unit.SetFile("file.template")
	unit.SetTarget("file.target")
	env := script.NewScriptEnvironmentWithGlobals(nil)
	err := unit.Render("doesnt-exist", env)
	c.Assert(err, IsNil)
}

func (s *testSuite) Test_Render_fails_with_invalid_mapping(c *C) {
	unit := NewTemplate()
	unit.SetFile("file.template")
	unit.SetTarget("file.target")
	unit.SetMapping(map[string]interface{}{
		"invalid": "$doesnt$work",
	})
	env := script.NewScriptEnvironmentWithGlobals(nil)
	err := unit.Render("build", env)
	c.Assert(err, Not(IsNil))
}

func (s *testSuite) Test_Render_Golden_Path(c *C) {
	os.RemoveAll("testdata/target.txt")
	unit := NewTemplate()
	unit.SetFile("testdata/helloworld.mustache")
	unit.SetTarget("testdata/target.txt")
	unit.SetScopes([]string{"myscope"})
	unit.SetMapping(map[string]interface{}{
		"who": "$this.variable",
	})
	thisDict := map[string]script.Script{
		"variable": script.LiftString("scripted world"),
	}
	globalsDict := map[string]script.Script{
		"this": script.LiftDict(thisDict),
	}
	env := script.NewScriptEnvironmentWithGlobals(globalsDict)
	err := unit.Render("myscope", env)
	c.Assert(err, IsNil)
	result, err := ioutil.ReadFile("testdata/target.txt")
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, "Hello scripted world\n")
	os.RemoveAll("testdata/target.txt")
}
