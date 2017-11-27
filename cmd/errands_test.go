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

package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/ankyra/escape/model/escape_plan"

	. "gopkg.in/check.v1"
)

func (s *suite) Test_List_Errands_Local(c *C) {
	cmd := RootCmd
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{"errands", "list", "--local"})

	c.Assert(ioutil.WriteFile("test.sh", []byte("#!/bin/bash\necho hello"), 0644), IsNil)
	plan := escape_plan.NewEscapePlan()
	plan.Name = "test"
	plan.Version = "1"
	plan.Errands["test-errand"] = map[string]interface{}{
		"script": "test.sh",
	}
	c.Assert(ioutil.WriteFile("escape.yml", plan.ToYaml(), 0644), IsNil)
	c.Assert(cmd.Execute(), IsNil)
	c.Assert(buf.String(), Equals, "")
	os.Remove("escape.yml")
	os.Remove("test.sh")
}

func (s *suite) Test_List_Errands_Local_no_errands(c *C) {
	cmd := RootCmd
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{"errands", "list", "--local"})

	plan := escape_plan.NewEscapePlan()
	plan.Name = "test"
	plan.Version = "1"
	c.Assert(ioutil.WriteFile("escape.yml", plan.ToYaml(), 0644), IsNil)
	c.Assert(cmd.Execute(), IsNil)
	c.Assert(buf.String(), Equals, "")
	os.Remove("escape.yml")
}

func (s *suite) Test_List_Errands_Local_missing_escape_plan(c *C) {
	cmd := RootCmd
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{"errands", "list", "--local"})
	c.Assert(cmd.Execute().Error(), Equals,
		"Escape plan 'escape.yml' was not found. Use 'escape plan init' to create it.")
}

func (s *suite) Test_List_Errands_missing_deployment(c *C) {
	cmd := RootCmd
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{"errands", "list", "-d", "test-deployment"})
	c.Assert(cmd.Execute().Error(), Equals, "The deployment 'test-deployment' could not be found in environment 'dev'")
}

func (s *suite) Test_Run_Errands_missing_deployment_name(c *C) {
	cmd := RootCmd
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{"errands", "run", "errand-name"})
	c.Assert(cmd.Execute().Error(), Equals, "Missing deployment name")
}
