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

package inventory

import (
	"github.com/ankyra/escape/model/inventory/local"
	"github.com/ankyra/escape/model/inventory/proxy"
	"github.com/ankyra/escape/model/inventory/remote"
	"github.com/ankyra/escape/model/inventory/types"
)

func NewLocalInventory(baseDir string) types.Inventory {
	return local.NewLocalInventory(baseDir)
}

func NewRemoteInventory(apiServer, authToken, basicAuthUsername, basicAuthPassword string, insecureSkipVerify bool) types.Inventory {
	return remote.NewRemoteInventory(apiServer, authToken, basicAuthUsername, basicAuthPassword, insecureSkipVerify)
}

func NewInventoryProxy(inv types.Inventory, proxiedNamespaces []string, proxyInv types.Inventory) types.Inventory {
	return proxy.NewInventoryProxy(inv, proxyInv, proxiedNamespaces)
}
