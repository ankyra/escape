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

package controllers

import (
	"github.com/ankyra/escape/model"
	. "github.com/ankyra/escape/model/interfaces"
)

type DepsController struct{}

func (DepsController) Fetch(context Context) error {
	context.PushLogRelease(context.GetEscapePlan().GetReleaseId())
	context.PushLogSection("Build")
	context.Log("fetch.start", nil)
	if len(context.GetEscapePlan().Depends) == 0 {
		return nil
	}
	depends, err := context.GetEscapePlan().GetDependencies()
	if err != nil {
		return err
	}
	err = model.DependencyResolver{}.Resolve(context.GetEscapeConfig(), depends)
	if err != nil {
		return err
	}
	context.Log("fetch.finished", nil)
	context.PopLogRelease()
	context.PopLogSection()
	return nil
}
