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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ankyra/escape/model/inventory"
	"github.com/ankyra/escape/util"
)

type EscapeConfig struct {
	ActiveProfile string                          `json:"current_profile"`
	Profiles      map[string]*EscapeProfileConfig `json:"profiles"`
	saveLocation  string                          `json:"-"`
}

type EscapeProfileConfig struct {
	ApiServer          string `json:"api_server"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	AuthToken          string `json:"escape_auth_token"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
	parent             *EscapeConfig
}

func NewEscapeConfig() *EscapeConfig {
	cfg := &EscapeConfig{}
	cfg.Profiles = map[string]*EscapeProfileConfig{
		"default": newEscapeProfileConfig(cfg),
	}
	cfg.ActiveProfile = "default"
	return cfg
}

func newEscapeProfileConfig(cfg *EscapeConfig) *EscapeProfileConfig {
	profile := &EscapeProfileConfig{
		ApiServer: os.Getenv("ESCAPE_API_SERVER"),
		Username:  os.Getenv("ESCAPE_USERNAME"),
		Password:  os.Getenv("ESCAPE_PASSWORD"),
		AuthToken: os.Getenv("ESCAPE_AUTH_TOKEN"),
		parent:    cfg,
	}
	if profile.ApiServer == "" {
		profile.ApiServer = "https://escape.ankyra.io"
	}
	return profile
}

func (c *EscapeConfig) GetInventory() inventory.Inventory {
	return c.GetCurrentProfile().GetInventory()
}

func (e *EscapeConfig) GetCurrentProfile() *EscapeProfileConfig {
	return e.Profiles[e.ActiveProfile]
}

func (e *EscapeConfig) LoadConfig(cfgFile, cfgProfile string) error {
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
		e.ActiveProfile = cfgProfile
	}
	if _, ok := e.Profiles[e.ActiveProfile]; !ok {
		return fmt.Errorf("Referenced profile '%s' was not found in the Escape configuration file.", e.ActiveProfile)
	}
	for _, t := range e.Profiles {
		t.parent = e
	}
	return nil
}

func (e *EscapeConfig) FromJson(cfgFile string) error {
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, e)
}

func (e *EscapeConfig) Save() error {
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

func (t *EscapeProfileConfig) ToJson() string {
	str, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (t *EscapeProfileConfig) GetInventory() inventory.Inventory {
	return inventory.NewRemoteInventory(t.ApiServer, t.AuthToken, t.InsecureSkipVerify)
}

func (t *EscapeProfileConfig) Save() error {
	return t.parent.Save()
}
func (t *EscapeProfileConfig) GetApiServer() string {
	return t.ApiServer
}
func (t *EscapeProfileConfig) GetUsername() string {
	return t.Username
}
func (t *EscapeProfileConfig) GetPassword() string {
	return t.Password
}
func (t *EscapeProfileConfig) GetAuthToken() string {
	return t.AuthToken
}
func (t *EscapeProfileConfig) SetApiServer(v string) {
	t.ApiServer = v
}
func (t *EscapeProfileConfig) SetUsername(v string) {
	t.Username = v
}
func (t *EscapeProfileConfig) SetPassword(v string) {
	t.Password = v
}
func (t *EscapeProfileConfig) SetAuthToken(v string) {
	t.AuthToken = v
}
func (t *EscapeProfileConfig) GetInsecureSkipVerify() bool {
	return t.InsecureSkipVerify
}
func (t *EscapeProfileConfig) SetInsecureSkipVerify(v bool) {
	t.InsecureSkipVerify = v
}
