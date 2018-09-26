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

package loggers

import "github.com/ankyra/escape/util/logger/api"

type LoggerDummy struct{}

func NewLoggerDummy() api.Logger {
	return &LoggerDummy{}
}

func (l *LoggerDummy) Log(key string, values map[string]string) {
}

func (l *LoggerDummy) PushSection(s string) {
}

func (l *LoggerDummy) Close() {
}

func (l *LoggerDummy) PopSection() {
}

func (l *LoggerDummy) PushRelease(s string) {
}

func (l *LoggerDummy) PopRelease() {
}

func (l *LoggerDummy) SetLogLevel(logLevel string) {
}
