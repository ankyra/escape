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
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/variable"
	"github.com/ankyra/escape-client/util"
	"os"
	"regexp"
	"strings"
)

type PackerReleaseType struct {
	ImageAlreadyExists *regexp.Regexp
	ImageCreated       *regexp.Regexp
}

func NewPackerReleaseType() ReleaseType {
	imageExistsReg := regexp.MustCompile(`Image (.*) already exists`)
	imageCreatedReg := regexp.MustCompile(`A disk image was created: (.*)\s*`)
	return &PackerReleaseType{
		ImageAlreadyExists: imageExistsReg,
		ImageCreated:       imageCreatedReg,
	}
}

func (a *PackerReleaseType) GetType() string {
	return "packer"
}

func (a *PackerReleaseType) InitEscapePlan(plan EscapePlan) {
	plan.SetPath(plan.GetBuild() + ".json")
}

func (a *PackerReleaseType) CompileMetadata(plan EscapePlan, metadata ReleaseMetadata) error {
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

func (a *PackerReleaseType) Run(ctx RunnerContext) (*map[string]interface{}, error) {
	outputs := map[string]interface{}{}
	env := a.buildEnvironment(ctx.GetReleaseMetadata(), *ctx.GetBuildInputs())
	stdout, status_code := a.runPacker(ctx.GetReleaseMetadata(), env, ctx.Logger())
	var image *string
	if status_code == 0 {
		image = a.readImageFromOutput(stdout)
	} else {
		image = a.recoverImageFromImageExistsFailure(stdout)
	}
	if image != nil {
		outputs["image"] = *image
	} else {
		return nil, fmt.Errorf("Packer build did not complete successfully")
	}
	return &outputs, nil
}

func (a *PackerReleaseType) Destroy(ctx RunnerContext) error {
	return nil
}

func (a *PackerReleaseType) buildEnvironment(metadata ReleaseMetadata, inputs map[string]interface{}) []string {
	env := os.Environ()
	for k, v := range inputs {
		switch v.(type) {
		case string:
			env = append(env, strings.ToUpper(k)+"="+v.(string))
		default:
			panic("yo, expecting a string fam")
		}
	}
	return env
}

func (a *PackerReleaseType) runPacker(metadata ReleaseMetadata, env []string, log Logger) (string, int) {
	cmd := []string{"packer", "build", metadata.GetPath()}
	p := util.NewProcessRecorder()
	stdout, err := p.Record(cmd, env, log)
	if err != nil {
		return stdout, 1
	}
	return stdout, 0
}
func (a *PackerReleaseType) recoverImageFromImageExistsFailure(stdout string) *string {
	images := a.ImageAlreadyExists.FindAllStringSubmatch(stdout, -1)
	if images == nil {
		return nil
	}
	return &images[0][1]
}
func (a *PackerReleaseType) readImageFromOutput(stdout string) *string {
	images := a.ImageCreated.FindAllStringSubmatch(stdout, -1)
	if images == nil {
		return nil
	}
	return &images[0][1]
}
