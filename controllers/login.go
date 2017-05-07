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
	"fmt"

	. "github.com/ankyra/escape-client/model/interfaces"
)

type LoginController struct{}

func (LoginController) Login(context Context, url, username, password string, storeCredentials bool) error {
	cfg := context.GetEscapeConfig().GetCurrentTarget()
	if username == "" {
		username = cfg.GetUsername()
	}
	if password == "" {
		password = cfg.GetPassword()
	}
	if username == "" {
		return fmt.Errorf("Missing username")
	}
	if password == "" {
		return fmt.Errorf("Missing password")
	}
	return context.GetClient().Login(url, username, password, storeCredentials)
}
