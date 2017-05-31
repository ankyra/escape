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

package compiler

import (
	"encoding/base64"
	"github.com/ankyra/escape-client/model/escape_plan"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type suite struct{}

var _ = Suite(&suite{})

func (s *suite) Test_Compile_Logo(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Logo = "testdata/logo.png"
	ctx := NewCompilerContext(plan, nil, "_")
	c.Assert(compileLogo(ctx), IsNil)
	data, err := ioutil.ReadFile(plan.Logo)
	c.Assert(err, IsNil)
	expected := base64.StdEncoding.EncodeToString(data)
	c.Assert(ctx.Metadata.Logo, Equals, expected)
}

func (s *suite) Test_Compile_Logo_does_nothing_if_no_logo_is_set(c *C) {
	plan := escape_plan.NewEscapePlan()
	ctx := NewCompilerContext(plan, nil, "_")
	c.Assert(compileLogo(ctx), IsNil)
}

func (s *suite) Test_Compile_Logo_fails_if_path_doesnt_exist(c *C) {
	plan := escape_plan.NewEscapePlan()
	plan.Logo = "testdata/doesnt_exist.png"
	ctx := NewCompilerContext(plan, nil, "_")
	c.Assert(compileLogo(ctx).Error(), Equals, "Referenced logo 'testdata/doesnt_exist.png' does not exist")
}

func (s *suite) Test_Compile_Logo_fails_if_cant_read_path(c *C) {
	path := "testdata/unreadable.png"
	plan := escape_plan.NewEscapePlan()
	plan.Logo = path
	c.Assert(os.Chmod(path, 0), IsNil)
	ctx := NewCompilerContext(plan, nil, "_")
	c.Assert(compileLogo(ctx).Error(), Equals, "Couldn't read logo 'testdata/unreadable.png': open testdata/unreadable.png: permission denied")
	c.Assert(os.Chmod(path, 0644), IsNil)
}
