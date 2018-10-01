package config

import (
	"encoding/json"
	"os"

	"github.com/ankyra/escape/model/inventory"
	"github.com/ankyra/escape/model/inventory/types"
	"github.com/ankyra/escape/model/paths"
)

type InventoryType string

var LocalInventory InventoryType = "local"
var RemoteInventory InventoryType = "remote"

type EscapeConfigProfile struct {
	InventoryType         InventoryType `json:"inventory_type"`
	ApiServer             string        `json:"api_server"`
	AuthToken             string        `json:"escape_auth_token"`
	BasicAuthUsername     string        `json:"basic_auth_username"`
	BasicAuthPassword     string        `json:"basic_auth_password"`
	InsecureSkipVerify    bool          `json:"insecure_skip_verify"`
	StatePath             string        `json:"state_path"`
	LocalInventoryBaseDir string        `json:"local_inventory_base_dir"`
	parent                *EscapeConfig
}

func newEscapeConfigProfile(cfg *EscapeConfig) *EscapeConfigProfile {
	profile := &EscapeConfigProfile{
		ApiServer:         os.Getenv("ESCAPE_API_SERVER"),
		AuthToken:         os.Getenv("ESCAPE_AUTH_TOKEN"),
		BasicAuthUsername: os.Getenv("BASIC_AUTH_USERNAME"),
		BasicAuthPassword: os.Getenv("BASIC_AUTH_PASSWORD"),
	}
	return profile.fix(cfg)
}

func (t *EscapeConfigProfile) fix(cfg *EscapeConfig) *EscapeConfigProfile {
	t.parent = cfg
	if t.InventoryType == "" {
		t.InventoryType = RemoteInventory
	}
	if t.ApiServer == "" {
		t.ApiServer = "https://escape.ankyra.io"
	}
	if t.StatePath == "" {
		t.StatePath = paths.NewPath().GetDefaultStateLocation()
	}
	return t
}

func (t *EscapeConfigProfile) ToJson() string {
	str, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (t *EscapeConfigProfile) GetInventory() types.Inventory {
	if t.InventoryType == LocalInventory {
		return inventory.NewLocalInventory(t.LocalInventoryBaseDir)
	}
	return inventory.NewRemoteInventory(t.ApiServer, t.AuthToken, t.BasicAuthUsername, t.BasicAuthPassword, t.InsecureSkipVerify)
}

func (t *EscapeConfigProfile) Save() error {
	return t.parent.Save()
}
func (t *EscapeConfigProfile) GetApiServer() string {
	return t.ApiServer
}
func (t *EscapeConfigProfile) GetAuthToken() string {
	return t.AuthToken
}
func (t *EscapeConfigProfile) SetApiServer(v string) {
	t.ApiServer = v
}
func (t *EscapeConfigProfile) SetAuthToken(v string) {
	t.AuthToken = v
}
func (t *EscapeConfigProfile) GetInsecureSkipVerify() bool {
	return t.InsecureSkipVerify
}
func (t *EscapeConfigProfile) SetInsecureSkipVerify(v bool) {
	t.InsecureSkipVerify = v
}
func (t *EscapeConfigProfile) SetBasicAuthCredentials(username, password string) {
	t.BasicAuthUsername = username
	t.BasicAuthPassword = password
}
func (t *EscapeConfigProfile) GetStatePath() string {
	return t.StatePath
}
