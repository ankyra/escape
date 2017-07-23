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
	RunningTestStep = "running_test_step"

	// Test and smoke failure
	TestFailure = "test_failure"

	// Destroy phase
	RunningPreDestroyStep  = "running_pre_destroy_step"
	RunningMainDestroyStep = "running_main_destroy_step"
	RunningPostDestroyStep = "running_post_destroy_step"

	// Destroy failure
	DestroyFailure = "destroy_failure"
)

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
	Code      StatusCode `json:"status"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`
	Data      string     `json:"data,omitempty"`
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
