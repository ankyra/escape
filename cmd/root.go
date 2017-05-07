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
	"os"

	"github.com/ankyra/escape-client/model"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/util"
	"github.com/spf13/cobra"
)

var cfgFile, cfgProfile, cfgLogLevel string
var context Context

var RootCmd = &cobra.Command{
	Use:           "escape-client",
	SilenceErrors: true,
	SilenceUsage:  true,
	Short:         "Package and deployment manager",
	Long: `Escape v` + util.EscapeVersion + ` 

Escape is a tool that can be used to version, package, build, release, 
deploy and operate software in the large and the small. Software of all sizes. 
Everyone welcome.

    ` + util.EscapeLogo + `
    
Website: http://escape.ankyra.io/
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		context = model.NewContext()
		err := context.LoadEscapeConfig(cfgFile, cfgProfile)
		if err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if context != nil {
			context.Log("error", map[string]string{
				"error": err.Error(),
			})
		} else {
			RootCmd.UsageFunc()(RootCmd)
		}
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "~/.escape_config", "Global Escape configuration file (default is $HOME/.escape-config)")
	RootCmd.PersistentFlags().StringVar(&cfgProfile, "profile", "default", "Configuration profile")
	RootCmd.PersistentFlags().StringVarP(&cfgLogLevel, "level", "l", "info", "Log level: debug, info, warn, error")
}
