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

/*
Downloading files at build or deployment time is one of those common tasks
that Escape tries to cover.

## Escape Plan

Downloads are configured in the Escape Plan under the
[`downloads`](/docs/reference/escape-plan/#downloads) field.

*/
type DownloadConfig struct {
	// The URL to download from. This field is required.
	//
	// Example: `https://www.google.com/`
	URL string `json:"url"`

	// The destination path.
	Dest string `json:"dest"`

	// Overwrite the destination path if it already exists.
	OverwriteExistingDest bool `json:"overwrite" yaml:"overwrite"`

	// Only perform this download if none of the paths in this list exist.
	// Supports glob patterns (for example: `"*.zip"`)
	IfNotExists []string `json:"if_not_exists" yaml:"if_not_exists"`

	// Should Escape try and unpack the destination path after download?
	// Supported extensions: `.zip`, `.tgz`, `.tar.gz`, `.tar`.
	Unpack bool `json:"unpack"`

	// Only perform this download if the platform matches this value.
	// Can be used to do platform dependent builds.
	Platform string `json:"platform"`

	// Only perform this download if the architecture matches this string.
	// Can be used to do architecture dependent builds.
	Arch string `json:"arch"`

	// A list of scopes (`build`, `deploy`) that defines during which stage(s)
	// this download should be performed.
	Scopes []string `json:"scopes" yaml:"scopes"`
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
