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
	"time"
)

type StatusCode string

const (
	// Initial state
	Empty StatusCode = "empty"

	// Configured, but not running
	Pending = "pending"

	// Build and deployment phase
	RunningPreStep  = "running_pre_step"
	RunningMainStep = "running_main_step"
	RunningPostStep = "running_post_step"

	// Build and deployment failure
	Failure = "failure"

	// Build, deployment, test and smoke success
	OK = "ok"

	// Test and smoke phase
	TestPending     = "test_pending"
	RunningTestStep = "running_test_step"

	// Test and smoke failure
	TestFailure = "test_failure"

	// Destroy phase
	DestroyPending          = "destroy_pending"
	DestroyAndDeletePending = "destroy_and_delete_pending"
	RunningPreDestroyStep   = "running_pre_destroy_step"
	RunningMainDestroyStep  = "running_main_destroy_step"
	RunningPostDestroyStep  = "running_post_destroy_step"

	// Destroy failure
	DestroyFailure = "destroy_failure"
)

// Can you go from one state to another?
var StatusTransitions = map[StatusCode][]StatusCode{
	Empty:          []StatusCode{RunningPreStep, Pending},
	Pending:        []StatusCode{RunningPreStep},
	OK:             []StatusCode{RunningPreStep, RunningTestStep, DestroyPending, DestroyAndDeletePending, TestPending, Pending, RunningPreDestroyStep},
	Failure:        []StatusCode{RunningPreStep, Pending, DestroyPending, DestroyAndDeletePending, RunningPreDestroyStep},
	TestFailure:    []StatusCode{RunningTestStep, RunningPreStep, DestroyPending, DestroyAndDeletePending, TestPending, Pending, RunningPreDestroyStep},
	DestroyFailure: []StatusCode{RunningPreDestroyStep, DestroyPending, DestroyAndDeletePending, Pending, RunningPreStep},

	RunningPreStep:  []StatusCode{RunningMainStep, Failure},
	RunningMainStep: []StatusCode{RunningPostStep, Failure},
	RunningPostStep: []StatusCode{RunningTestStep, OK, Failure},

	RunningTestStep: []StatusCode{OK, TestFailure},

	DestroyPending:          []StatusCode{RunningPreDestroyStep},
	DestroyAndDeletePending: []StatusCode{RunningPreDestroyStep},
	RunningPreDestroyStep:   []StatusCode{RunningMainDestroyStep, DestroyFailure},
	RunningMainDestroyStep:  []StatusCode{RunningPostDestroyStep, DestroyFailure},
	RunningPostDestroyStep:  []StatusCode{Empty, DestroyFailure},
}

// Can you go from s1 -> s2?
func StatusTransitionAllowed(s1, s2 StatusCode) bool {
	transitions := StatusTransitions[s1]
	for _, allowed := range transitions {
		if allowed == s2 {
			return true
		}
	}
	return false
}

var ErrorStatuses = map[StatusCode]bool{
	Failure:        true,
	TestFailure:    true,
	DestroyFailure: true,
}
var OKStatuses = map[StatusCode]bool{
	Empty:   true,
	Pending: true,
	OK:      true,
}
var RunningStatus = map[StatusCode]bool{}

type Status struct {
	Code       StatusCode `json:"status"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
	UpdatedBy  string     `json:"updated_by,omitempty"`
	Data       string     `json:"data,omitempty"`
	TryAgainAt *time.Time `json:"try_again_at,omitempty"`
	Tried      int        `json:"tried,omitempty"`
}

func NewStatus(code StatusCode) *Status {
	return &Status{
		Code:      code,
		UpdatedAt: time.Now(),
	}
}

func (s *Status) IsError() bool {
	_, found := ErrorStatuses[s.Code]
	return found
}

func (s *Status) IsOK() bool {
	_, found := OKStatuses[s.Code]
	return found
}

func (s *Status) IsRunning() bool {
	return !(s.IsError() || s.IsOK())
}

func (s *Status) IsOneOf(codes ...StatusCode) bool {
	for _, c := range codes {
		if s.Code == c {
			return true
		}
	}
	return false
}
