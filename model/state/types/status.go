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

package types

import (
	"time"
)

type StatusCode int

const (
	// Initial state
	Empty StatusCode = iota

	// Configured, but not running
	Pending

	// Build and deployment phase
	RunningPreStep
	RunningMainStep
	RunningPostStep

	// Build and deployment failure
	Failure

	// Build, deployment, test and smoke success
	OK

	// Test and smoke phase
	RunningTestStep

	// Test and smoke failure
	TestFailure

	// Destroy phase
	RunningPreDestroyStep
	RunningMainDestroyStep
	RunningPostDestroyStep

	// Destroy failure
	DestroyFailure
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
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
	Data      string     `json:"data"`
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
