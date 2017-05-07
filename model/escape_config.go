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

package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/registry"
	"github.com/ankyra/escape-client/util"
)

type escapeConfig struct {
	ActiveTarget string                         `json:"current_target"`
	Targets      map[string]*escapeTargetConfig `json:"targets"`
	Context      Context                        `json:"-"`
	saveLocation string                         `json:"-"`
	Client       Client                         `json:"-"`
}

type escapeTargetConfig struct {
	ApiServer      string  `json:"api_server"`
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	AuthToken      string  `json:"escape_auth_token"`
	StorageBackend string  `json:"storage_backend"`
	GcsBucketUrl   string  `json:"bucket_url"`
	Context        Context `json:"-"`
}

func NewEscapeConfig(context Context) EscapeConfig {
	cfg := &escapeConfig{
		Context: context,
		Targets: map[string]*escapeTargetConfig{
			"default": newEscapeTargetConfig(context),
		},
		ActiveTarget: "default",
	}
	cfg.Client = registry.NewClient(cfg)
	return cfg
}

func newEscapeTargetConfig(context Context) *escapeTargetConfig {
	target := &escapeTargetConfig{
		Context:        context,
		ApiServer:      os.Getenv("ESCAPE_API_SERVER"),
		Username:       os.Getenv("ESCAPE_USERNAME"),
		Password:       os.Getenv("ESCAPE_PASSWORD"),
		AuthToken:      os.Getenv("ESCAPE_AUTH_TOKEN"),
		StorageBackend: os.Getenv("ESCAPE_STORAGE_BACKEND"),
		GcsBucketUrl:   os.Getenv("ESCAPE_BUCKET_URL"),
	}
	return target
}

func (c *escapeConfig) GetClient() Client {
	return c.Client
}

func (e *escapeConfig) GetCurrentTarget() EscapeTargetConfig {
	return e.Targets[e.ActiveTarget]
}

func (e *escapeConfig) LoadConfig(cfgFile, cfgProfile string) error {
	if len(cfgFile) > 2 && cfgFile[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		cfgFile = strings.Replace(cfgFile, "~/", dir+"/", 1)
	}
	cfgFile, err := filepath.Abs(cfgFile)
	if err != nil {
		return err
	}
	e.saveLocation = cfgFile
	if util.PathExists(cfgFile) {
		err := e.FromJson(cfgFile)
		if err != nil {
			return fmt.Errorf("Couldn't parse Escape configuration file '%s': %s", cfgFile, err.Error())
		}
	}
	if cfgProfile != "" {
		e.ActiveTarget = cfgProfile
	}
	if _, ok := e.Targets[e.ActiveTarget]; !ok {
		return fmt.Errorf("Referenced target '%s' was not found in the Escape configuration file.", e.ActiveTarget)
	}
	for _, t := range e.Targets {
		t.Context = e.Context
	}
	return nil
}

func (e *escapeConfig) FromJson(cfgFile string) error {
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, e)
}

func (e *escapeConfig) Save() error {
	if e.saveLocation == "" {
		return fmt.Errorf("Save location has not been set")
	}
	str, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not convert escape config to json: %s", err.Error())
	}
	mode := os.FileMode(0600)
	st, err := os.Stat(e.saveLocation)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("Could not stat Escape config file at '%s': %s", e.saveLocation, err.Error())
		}
	} else {
		mode = st.Mode()
	}
	return ioutil.WriteFile(e.saveLocation, str, mode)
}

func (t *escapeTargetConfig) ToJson() string {
	str, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (t *escapeTargetConfig) Save() error {
	return t.Context.GetEscapeConfig().Save()
}
func (t *escapeTargetConfig) GetApiServer() string {
	return t.ApiServer
}
func (t *escapeTargetConfig) GetUsername() string {
	return t.Username
}
func (t *escapeTargetConfig) GetPassword() string {
	return t.Password
}
func (t *escapeTargetConfig) GetAuthToken() string {
	return t.AuthToken
}
func (t *escapeTargetConfig) GetStorageBackend() string {
	return t.StorageBackend
}
func (t *escapeTargetConfig) GetGcsBucketUrl() string {
	return t.GcsBucketUrl
}

func (t *escapeTargetConfig) SetApiServer(v string) {
	t.ApiServer = v
}
func (t *escapeTargetConfig) SetUsername(v string) {
	t.Username = v
}
func (t *escapeTargetConfig) SetPassword(v string) {
	t.Password = v
}
func (t *escapeTargetConfig) SetAuthToken(v string) {
	t.AuthToken = v
}
func (t *escapeTargetConfig) SetStorageBackend(v string) {
	t.StorageBackend = v
}
func (t *escapeTargetConfig) SetGcsBucketUrl(v string) {
	t.GcsBucketUrl = v
}
