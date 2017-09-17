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
	URL    string `json:"url"`
	Dest   string `json:"dest"`
	Unpack bool   `json:"unpack"`
}

func NewDownloadConfig(url string) *DownloadConfig {
	return &DownloadConfig{
		URL: url,
	}
}

func (d *DownloadConfig) ValidateAndFix() error {
	if d.URL == "" {
		return fmt.Errorf("Missing URL in download config.")
	}
	parsed, err := url.Parse(d.URL)
	if err != nil {
		return fmt.Errorf("Invalid URL '%s' in download config: %s", d.URL, err.Error())
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("Missing URL scheme in download config '%s'", d.URL)
	}
	return nil
}
