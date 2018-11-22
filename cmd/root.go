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

package cmd

import (
	"fmt"
	"os"

	"github.com/ankyra/escape/model"
	"github.com/ankyra/escape/util"
	"github.com/ankyra/escape/util/logger"
	"github.com/spf13/cobra"
)

var cfgFile, cfgProfile, cfgLogLevel, cfgLogger string
var cfgLogCollapse, jsonFlag bool
var context *model.Context

var RootCmd = &cobra.Command{
	Use:           "escape",
	SilenceErrors: true,
	SilenceUsage:  true,
	Short:         "Package and deployment manager",
	Long: `Escape v` + util.EscapeVersion + ` 

Escape is a tool to help with the release engineering, life-cycle management
and Continuous Delivery of software platforms and their artefacts.
    
See the documentation at https://escape.ankyra.io/docs/
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		context = model.NewContext()
		err := context.LoadEscapeConfig(cfgFile, cfgProfile)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		logger, err := logger.GetLogger(cfgLogger, cfgLogLevel, cfgLogCollapse)
		if err != nil {
			return err
		}
		context.SetLogger(logger)
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if context != nil {
			context.Log("error", map[string]string{
				"error": err.Error(),
			})
			context.Logger.Close()
		} else {
			RootCmd.UsageFunc()(RootCmd)
		}
		os.Exit(1)
	}
	if context != nil {
		context.Logger.Close()
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "~/.escape_config", "Local of the global Escape configuration file")
	RootCmd.PersistentFlags().StringVar(&cfgProfile, "profile", "", "Configuration profile")
	RootCmd.PersistentFlags().StringVarP(&cfgLogLevel, "level", "l", "info", "Log level: debug, success, info, warn, error")
	RootCmd.PersistentFlags().StringVarP(&cfgLogger, "logger", "", "default", "Logger: default, json")
	RootCmd.PersistentFlags().BoolVarP(&cfgLogCollapse, "collapse-logs", "", true, "Collapse log sections.")
}
