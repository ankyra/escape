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

package state

import (
	. "gopkg.in/check.v1"
	"time"
)

func (s *suite) Test_Status_Sets_Time(c *C) {
	now := time.Now()
	st := NewStatus(OK)
	c.Assert(st.UpdatedAt.Sub(now) > 1, Equals, true)
}

func (s *suite) Test_Status_IsOneOf(c *C) {
	c.Assert(NewStatus(OK).IsOneOf(), Equals, false)
	c.Assert(NewStatus(OK).IsOneOf(Empty), Equals, false)
	c.Assert(NewStatus(OK).IsOneOf(Empty, Pending), Equals, false)
	c.Assert(NewStatus(OK).IsOneOf(Empty, Pending, OK), Equals, true)
	c.Assert(NewStatus(OK).IsOneOf(OK), Equals, true)

}

func (s *suite) Test_Status_IsError(c *C) {
	trueCases := []StatusCode{
		Failure,
		TestFailure,
		DestroyFailure,
	}
	falseCases := []StatusCode{
		Empty,
		Pending,
		RunningPreStep,
		RunningMainStep,
		RunningPostStep,
		OK,
		RunningTestStep,
		RunningPreDestroyStep,
		RunningMainDestroyStep,
		RunningPostDestroyStep,
	}
	for _, t := range trueCases {
		st := NewStatus(t)
		c.Assert(st.IsError(), Equals, true)
	}
	for _, t := range falseCases {
		st := NewStatus(t)
		c.Assert(st.IsError(), Equals, false)
	}
}

func (s *suite) Test_Status_IsOK(c *C) {
	trueCases := []StatusCode{
		Empty,
		Pending,
		OK,
	}
	falseCases := []StatusCode{
		Failure,
		TestFailure,
		DestroyFailure,
		RunningPreStep,
		RunningMainStep,
		RunningPostStep,
		RunningTestStep,
		RunningPreDestroyStep,
		RunningMainDestroyStep,
		RunningPostDestroyStep,
	}
	for _, t := range trueCases {
		st := NewStatus(t)
		c.Assert(st.IsOK(), Equals, true)
	}
	for _, t := range falseCases {
		st := NewStatus(t)
		c.Assert(st.IsOK(), Equals, false)
	}
}

func (s *suite) Test_Status_IsRunning(c *C) {
	trueCases := []StatusCode{
		RunningPreStep,
		RunningMainStep,
		RunningPostStep,
		RunningTestStep,
		RunningPreDestroyStep,
		RunningMainDestroyStep,
		RunningPostDestroyStep,
	}
	falseCases := []StatusCode{
		Empty,
		Pending,
		OK,
		Failure,
		TestFailure,
		DestroyFailure,
	}
	for _, t := range trueCases {
		st := NewStatus(t)
		c.Assert(st.IsRunning(), Equals, true)
	}
	for _, t := range falseCases {
		st := NewStatus(t)
		c.Assert(st.IsRunning(), Equals, false)
	}
}
