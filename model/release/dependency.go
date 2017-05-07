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

package release

import (
	"errors"
	"strings"

	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/parsers"
)

type dependency struct {
	Type         string
	Build        string
	Version      string
	VariableName string
}

func NewDependencyFromMetadata(metadata ReleaseMetadata) Dependency {
	return &dependency{
		Type:    metadata.GetType(),
		Build:   metadata.GetName(),
		Version: metadata.GetVersion(),
	}
}

func NewDependencyFromString(str string) (Dependency, error) {
	parsed, err := parsers.ParseDependency(str)
	if err != nil {
		return nil, err
	}
	return &dependency{
		Type:         parsed.Type,
		Build:        parsed.Build,
		Version:      parsed.Version,
		VariableName: parsed.VariableName,
	}, nil
}

func (d *dependency) ResolveVersion(context Context) error {
	if d.Version != "latest" && !strings.HasSuffix(d.Version, "@") {
		return nil
	}
	backend := context.GetEscapeConfig().GetCurrentTarget().GetStorageBackend()
	if backend == "escape" {
		metadata, err := context.GetClient().ReleaseQuery(d.GetReleaseId())
		if err != nil {
			return err
		}
		d.Version = metadata.GetVersion()
	} else if backend == "" || backend == "local" {
		return errors.New("Backend not implemented: " + backend)
	} else {
		return errors.New("Unsupported Escape storage backend: " + backend)
	}
	return nil
}

func (d *dependency) GetType() string {
	return d.Type
}
func (d *dependency) GetBuild() string {
	return d.Build
}
func (d *dependency) GetVariableName() string {
	return d.VariableName
}
func (d *dependency) GetVersion() string {
	return d.Version
}

func (d *dependency) GetReleaseId() string {
	version := "v" + d.Version
	if d.Version == "latest" {
		version = d.Version
	}
	return d.Type + "-" + d.Build + "-" + version
}
func (d *dependency) GetVersionlessReleaseId() string {
	return d.Type + "-" + d.Build
}
