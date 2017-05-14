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
	"fmt"
	"strings"

	"github.com/ankyra/escape-client/model/escape_plan"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/variable"
	"github.com/ankyra/escape-client/util"
)

type DockerReleaseType struct{}

func (a *DockerReleaseType) GetType() string {
	return "docker"
}

func (a *DockerReleaseType) InitEscapePlan(plan *escape_plan.EscapePlan) {
	plan.SetIncludes([]string{"Dockerfile"})
}

func checkExistingVariable(v *variable.Variable, id string, typ string) (bool, error) {
	if v.GetId() == id {
		if v.GetType() != typ {
			return false, fmt.Errorf("Declared variable '%s' is not of expected type '%s', but '%s'", v.GetId(), typ, v.GetType())
		}
		return true, nil
	}
	return false, nil
}

func (a *DockerReleaseType) CompileMetadata(plan *escape_plan.EscapePlan, metadata ReleaseMetadata) error {
	dockerRepoFound := false
	dockerCmdFound := false
	for _, i := range metadata.GetInputs() {
		found, err := checkExistingVariable(i, "docker_repository", "string")
		if err != nil {
			return err
		}
		dockerRepoFound = dockerRepoFound || found
		found, err = checkExistingVariable(i, "docker_cmd", "list")
		if err != nil {
			return err
		}
		dockerCmdFound = dockerCmdFound || found
	}
	if !dockerRepoFound {
		v := variable.NewVariableFromString("docker_repository", "string")
		defaultValue := ""
		v.SetDefault(&defaultValue)
		metadata.AddInputVariable(v)
	}
	// TODO
	if !dockerCmdFound {
		v := variable.NewVariableFromString("docker_cmd", "list")
		defaultValue := []interface{}{}
		v.SetDefault(defaultValue)
		v.SetDescription("Overrides the default Docker command.")
		metadata.AddInputVariable(v)
	}
	for _, i := range metadata.GetOutputs() {
		found, err := checkExistingVariable(i, "image", "string")
		if err != nil {
			return err
		}
		if found {
			return nil
		}
	}
	v := variable.NewVariableFromString("image", "string")
	metadata.AddOutputVariable(v)
	return nil
}

func (a *DockerReleaseType) Run(ctx RunnerContext) (*map[string]interface{}, error) {
	metadata := ctx.GetReleaseMetadata()
	shouldPush := false
	repo, repoFound := (*ctx.GetBuildInputs())["docker_repository"]
	if repoFound && repo.(string) != "" {
		shouldPush = true
	}
	baseImageName := metadata.GetName()
	if shouldPush {
		baseImageName = repo.(string) + "/" + baseImageName
	}
	image := baseImageName + ":v" + metadata.GetVersion()
	log := ctx.Logger()
	docker, err := newDockerCmdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//        applog("build.docker", image=image, release=release_id)
	if err := docker.run([]string{"build", "-t", image, "."}, log); err != nil {
		return nil, err
	}
	if err := docker.run([]string{"tag", image, baseImageName + ":latest"}, log); err != nil {
		if err := docker.run([]string{"tag", "-f", image, baseImageName + ":latest"}, log); err != nil { // fallback for Docker <1.10
			return nil, err
		}
	}
	if shouldPush {
		if err := docker.run([]string{"push", image}, log); err != nil {
			return nil, err
		}
		if err := docker.run([]string{"push", baseImageName + ":latest"}, log); err != nil {
			return nil, err
		}
	} else {
		//            applog("build.docker_no_repository", release=release_id)
	}
	result := map[string]interface{}{"image": image}
	return &result, nil
}

func (a *DockerReleaseType) Destroy(ctx RunnerContext) error {
	outputs := ctx.GetBuildOutputs()
	if outputs == nil || (*outputs)["image"] == nil {
		return fmt.Errorf("Missing 'image' output variable.")
	}
	image, ok := (*outputs)["image"].(string)
	if !ok {
		return fmt.Errorf("Expecting 'image' output variable to be a string, but got %T.", (*outputs)["image"])
	}
	//            applog("destroy.docker_start", image=image, release=release_id)
	docker, err := newDockerCmdFromContext(ctx)
	if err != nil {
		return err
	}
	stdout, err := docker.record([]string{"images", "-q", image}, ctx.Logger())
	if err != nil {
		return err
	}
	stdout = strings.TrimSpace(stdout)
	if stdout == "" {
		//                applog("destroy.docker_already_removed", image=image, release=release_id)
		return nil
	}
	parts := strings.Split(image, ":")
	latest := strings.Join(parts[:len(parts)-1], ":") + ":latest"
	latestOut, err := docker.record([]string{"images", "-q", latest}, ctx.Logger())
	if err != nil {
		return err
	}
	latestOut = strings.TrimSpace(latestOut)
	isLatest := latestOut == stdout
	if err := docker.run([]string{"rmi", image}, ctx.Logger()); err != nil {
		return err
	}
	if isLatest {
		if err := docker.run([]string{"rmi", latest}, ctx.Logger()); err != nil {
			return err
		}
	}
	//            applog("destroy.docker_finished", image=image, release=release_id)
	return nil
}

func newDockerCmdFromContext(ctx RunnerContext) (*dockerCmd, error) {
	dockerCmd_, dockerCmdFound := (*ctx.GetBuildInputs())["docker_cmd"]
	if !dockerCmdFound {
		dockerCmd_ = []interface{}{"docker"}
	}
	dockerCmd := []string{}
	for _, s := range dockerCmd_.([]interface{}) {
		str, ok := s.(string)
		if !ok {
			return nil, fmt.Errorf("Expecting string in docker_cmd variable, got: %v [%T]", s, s)
		}
		dockerCmd = append(dockerCmd, str)
	}
	if len(dockerCmd) == 0 {
		dockerCmd = []string{"docker"}
	}
	return newDockerCmd(dockerCmd, ctx.GetPath().GetBaseDir()), nil
}

func newDockerCmd(cmd []string, workingDir string) *dockerCmd {
	return &dockerCmd{
		cmd:        cmd,
		workingDir: workingDir,
	}
}

type dockerCmd struct {
	cmd        []string
	workingDir string
}

func (d *dockerCmd) record(args []string, log util.Logger) (string, error) {
	cmd := []string{}
	for _, c := range d.cmd {
		cmd = append(cmd, c)
	}
	for _, a := range args {
		cmd = append(cmd, a)
	}
	proc := util.NewProcessRecorder()
	proc.SetWorkingDirectory(d.workingDir)
	return proc.Record(cmd, nil, log)
}

func (d *dockerCmd) run(args []string, log util.Logger) error {
	_, err := d.record(args, log)
	return err
}
