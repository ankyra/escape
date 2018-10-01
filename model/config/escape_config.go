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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ankyra/escape/model/inventory/types"
	"github.com/ankyra/escape/util"
)

type EscapeConfig struct {
	ActiveProfile string                          `json:"current_profile"`
	Profiles      map[string]*EscapeConfigProfile `json:"profiles"`
	saveLocation  string                          `json:"-"`
}

func NewEscapeConfig() *EscapeConfig {
	cfg := &EscapeConfig{}
	cfg.Profiles = map[string]*EscapeConfigProfile{
		"default": newEscapeConfigProfile(cfg),
	}
	cfg.ActiveProfile = "default"
	return cfg
}

func (c *EscapeConfig) NewProfile(profileName string) error {
	if c.Profiles[profileName] != nil {
		return fmt.Errorf("Profile already exists")
	}
	c.Profiles[profileName] = newEscapeConfigProfile(c)
	return nil
}

func (c *EscapeConfig) GetInventory() types.Inventory {
	return c.GetCurrentProfile().GetInventory()
}

func (e *EscapeConfig) GetCurrentProfile() *EscapeConfigProfile {
	return e.Profiles[e.ActiveProfile]
}

func (e *EscapeConfig) SetActiveProfile(profile string) error {
	if profile == "" {
		return nil
	}

	if e.Profiles[profile] == nil {
		return fmt.Errorf("Referenced profile '%s' was not found in the Escape configuration file.", profile)
	}
	e.ActiveProfile = profile
	return nil
}

func (e *EscapeConfig) LoadConfig(cfgFile string) error {
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
	} else {
		err := e.Save()
		if err != nil {
			return err
		}
	}
	if _, ok := e.Profiles[e.ActiveProfile]; !ok {
		return fmt.Errorf("Referenced profile '%s' was not found in the Escape configuration file.", e.ActiveProfile)
	}
	for _, t := range e.Profiles {
		t.fix(e)
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
