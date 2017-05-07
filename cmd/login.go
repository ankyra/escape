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

package cmd

import (
	"fmt"

	"github.com/ankyra/escape-client/controllers"
	"github.com/spf13/cobra"
)

var username, password, url string
var storeCredentials bool

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with an Escape server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if url == "" {
			return fmt.Errorf("Missing Escape server URL")
		}
		return controllers.LoginController{}.Login(context, url, username, password, storeCredentials)
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "The username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "The password")
	loginCmd.Flags().StringVarP(&url, "url", "e", "", "The Escape server URL")
	loginCmd.Flags().BoolVarP(&storeCredentials, "store", "s", false, "Store username and password in the Escape configuration file")
}
