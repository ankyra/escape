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
	"github.com/ankyra/escape-client/model/release"
	"github.com/ankyra/escape-client/util"
	"io/ioutil"
	"strings"
)

type KubeSpecReleaseType struct{}

func (a *KubeSpecReleaseType) GetType() string {
	return "kubespec"
}

func (a *KubeSpecReleaseType) InitEscapePlan(plan EscapePlan) {
	plan.SetIncludes([]string{"*.yml"})
	plan.SetConsumes([]string{"kubernetes"})
}

func (a *KubeSpecReleaseType) CompileMetadata(plan EscapePlan, metadata ReleaseMetadata) error {
	kubernetesConsumerFound := false
	credentialInputFound := false
	for _, c := range metadata.GetConsumes() {
		if c == "kubernetes" {
			kubernetesConsumerFound = true
		}
	}
	for _, i := range metadata.GetInputs() {
		found, err := checkExistingVariable(i, "account_credentials", "string")
		if err != nil {
			return err
		}
		credentialInputFound = credentialInputFound || found
	}
	if !kubernetesConsumerFound {
		consumes := metadata.GetConsumes()
		consumes = append(consumes, "kubernetes")
		metadata.SetConsumes(consumes)
	}
	if !credentialInputFound {
		v := release.NewVariableFromString("account_credentials", "string")
		defaultValue := "$kubernetes.outputs.kubernetes_kubectl.file"
		v.SetDefault(&defaultValue)
		v.SetDescription("The Kubernetes kubectl file.")
		v.SetSensitive(true)
		v.SetVisible(false)
		metadata.AddInputVariable(v)
	}
	return nil
}

func (a *KubeSpecReleaseType) Destroy(ctx RunnerContext) error {
	return nil
}

func (a *KubeSpecReleaseType) Run(ctx RunnerContext) (*map[string]interface{}, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}
	creds_, found := (*ctx.GetBuildInputs())["account_credentials"]
	if !found {
		return nil, fmt.Errorf("Missing 'account_credentials' input variable")
	}
	kubecfgPath := creds_.(string)
	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			continue
		} else if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			if name != "escape.yml" {
				proc := util.NewProcessRecorder()
				cmd := []string{"kubectl", "--kubeconfig", kubecfgPath,
					"apply", "-f", name}
				if err := proc.Run(cmd, nil, ctx.Logger()); err != nil {
					return nil, fmt.Errorf("While running kubectl on spec %s: %s", name, err.Error())
				}
			}
		}
	}
	return nil, nil
}
