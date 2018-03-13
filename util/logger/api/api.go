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

package api

type LogLevel int

const (
	DEBUG   = iota
	INFO    = iota
	SUCCESS = iota
	WARN    = iota
	ERROR   = iota
)

func StringToLogLevel(logLevel string) LogLevel {
	if logLevel == "debug" {
		return DEBUG
	} else if logLevel == "warn" {
		return WARN
	} else if logLevel == "success" {
		return SUCCESS
	} else if logLevel == "error" {
		return ERROR
	}
	return INFO
}

type LogEntry struct {
	Message      string
	SectionStack []string
	Release      string
	Collapse     bool
	LogLevel     LogLevel
}

type Logger interface {
	Log(string, map[string]string)
	PushSection(string)
	PopSection()
	PushRelease(string)
	PopRelease()
	SetLogLevel(string)
}

type LogConsumer interface {
	Consume(*LogEntry) error
}
