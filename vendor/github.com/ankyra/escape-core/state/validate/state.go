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

package validate

import (
	"fmt"
	"regexp"
)

func InvalidProjectNameError(name string) error {
	return fmt.Errorf("The project name '%s' is not allowed. Expected a string matching /%s/",
		name, projectNameRegexFmt)
}

func InvalidDeploymentNameError(name string) error {
	return fmt.Errorf("The deployment name '%s' is not allowed. Expected a string matching /%s/",
		name, deploymentNameRegexFmt)
}

func InvalidEnvironmentNameError(name string) error {
	return fmt.Errorf("The environment name '%s' is not allowed. Expected a string matching /%s/",
		name, environmentNameRegexFmt)
}

func InvalidStageNameError(name string) error {
	return fmt.Errorf("Invalid stage name '%s'. Expecting 'build' or 'deploy'.",
		name)
}

var projectNameRegexFmt = "^[a-zA-Z0-9]+[a-zA-Z0-9-_]+$"
var projectNameRegex = regexp.MustCompile(projectNameRegexFmt)

var environmentNameRegexFmt = "^[a-z]+[a-z0-9-_]*$"
var environmentNameRegex = regexp.MustCompile(environmentNameRegexFmt)

var deploymentNameRegexFmt = "^[a-zA-Z_]+[a-zA-Z0-9-_/]*$"
var deploymentNameRegex = regexp.MustCompile(deploymentNameRegexFmt)

func IsValidProjectName(name string) bool {
	return projectNameRegex.MatchString(name)
}

func IsValidDeploymentName(name string) bool {
	return deploymentNameRegex.MatchString(name)
}

func IsValidEnvironmentName(name string) bool {
	return environmentNameRegex.MatchString(name)
}

func IsValidStageName(name string) bool {
	return name == "build" || name == "deploy"
}
