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

package state

import (
	"github.com/ankyra/escape-core/state/validate"
	. "gopkg.in/check.v1"
)

func (s *suite) Test_newStage(c *C) {
	st := newStage()
	c.Assert(st.UserInputs, Not(IsNil))
	c.Assert(st.Inputs, Not(IsNil))
	c.Assert(st.Outputs, Not(IsNil))
	c.Assert(st.Providers, Not(IsNil))
	c.Assert(st.Status, Not(IsNil))
	c.Assert(st.Status.Code, DeepEquals, Empty)
}

func (s *suite) Test_Stage_validateAndFix_errors_on_invalid_name(c *C) {
	cases := []string{
		"",
		"aweoijaweioj",
		"  build",
		"  deploy  ",
		"deploy ",
		".",
	}
	for _, test := range cases {
		st := newStage()
		c.Assert(st.validateAndFix(test, nil, nil), DeepEquals, validate.InvalidStageNameError(test))
	}
}

func (s *suite) Test_Stage_validateAndFix_fixes_nils(c *C) {
	st := newStage()
	st.Name = DeployStage
	st.UserInputs = nil
	st.Inputs = nil
	st.Outputs = nil
	st.Providers = nil
	st.Deployments = nil
	st.Status = nil
	c.Assert(st.validateAndFix(DeployStage, nil, nil), IsNil)
	c.Assert(st.UserInputs, Not(IsNil))
	c.Assert(st.Inputs, Not(IsNil))
	c.Assert(st.Outputs, Not(IsNil))
	c.Assert(st.Providers, Not(IsNil))
	c.Assert(st.Status, Not(IsNil))
	c.Assert(st.Status.Code, DeepEquals, Empty)
}

func (s *suite) Test_Stage_validateAndFix_sets_Status_to_OK_if_nil_but_version_is_set(c *C) {
	st := newStage()
	st.Status = nil
	st.SetVersion("1.0")
	c.Assert(st.validateAndFix(DeployStage, nil, nil), IsNil)
	c.Assert(st.Status.Code, DeepEquals, StatusCode(OK))
}

func (s *suite) Test_Stage_SetInputs_inits_new_map_if_set_to_nil(c *C) {
	st := newStage()
	st.SetInputs(nil)
	c.Assert(st.Inputs, Not(IsNil))
}

func (s *suite) Test_Stage_SetUserInputs_inits_new_map_if_set_to_nil(c *C) {
	st := newStage()
	st.SetUserInputs(nil)
	c.Assert(st.UserInputs, Not(IsNil))
}

func (s *suite) Test_Stage_SetOutputs_inits_new_map_if_set_to_nil(c *C) {
	st := newStage()
	st.SetOutputs(nil)
	c.Assert(st.Outputs, Not(IsNil))
}
