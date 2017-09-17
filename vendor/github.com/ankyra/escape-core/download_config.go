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

package core

import (
	"fmt"
	"net/url"
)

type DownloadConfig struct {
	URL                   string   `json:"url"`
	Dest                  string   `json:"dest"`
	OverwriteExistingDest bool     `json:"overwrite" yaml:"overwrite"`
	IfNotExists           []string `json:"if_not_exists" yaml:"if_not_exists"`
	Unpack                bool     `json:"unpack"`
	Platform              string   `json:"platform"`
	Arch                  string   `json:"arch"`
	Scopes                []string `json:"scopes" yaml:"scopes"`
}

func NewDownloadConfig(url string) *DownloadConfig {
	return &DownloadConfig{
		URL:    url,
		Scopes: []string{"build", "deploy"},
	}
}

func (d *DownloadConfig) ValidateAndFix() error {
	if d.URL == "" {
		return fmt.Errorf("Missing URL in download config.")
	}
	if d.Dest == "" {
		return fmt.Errorf("Missing 'dest' in download config for '%s'", d.URL)
	}
	parsed, err := url.Parse(d.URL)
	if err != nil {
		return fmt.Errorf("Invalid URL '%s' in download config: %s", d.URL, err.Error())
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("Missing URL scheme in download config '%s'", d.URL)
	}
	if d.Scopes == nil || len(d.Scopes) == 0 {
		d.Scopes = []string{"build", "deploy"}
	}
	return nil
}

func (d *DownloadConfig) InScope(scope string) bool {
	for _, s := range d.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}
