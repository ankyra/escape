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

package logger

import (
	"fmt"

	"github.com/ankyra/escape/util/logger/api"
	"github.com/ankyra/escape/util/logger/consumers"
	"github.com/ankyra/escape/util/logger/loggers"
)

func GetLogger(logger string, logLevel string, collapse bool) (api.Logger, error) {
	var consumer api.LogConsumer
	if logger == "default" {
		consumer = consumers.NewFancyTerminalOutputLogConsumer(collapse)
	} else if logger == "json" {
		consumer = consumers.NewJSONLogConsumer()
	} else {
		return nil, fmt.Errorf("Unknown logger type '%s', expecting 'default' or 'json'.", logger)
	}
	result := loggers.NewLogger([]api.LogConsumer{consumer})
	result.SetLogLevel(logLevel)
	return result, nil
}
