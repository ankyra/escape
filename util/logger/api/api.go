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

import "time"

type LogLevel int

const (
	DEBUG   = iota
	INFO    = iota
	SUCCESS = iota
	WARN    = iota
	ERROR   = iota
)

func (l LogLevel) String() string {
	if l == DEBUG {
		return "debug"
	} else if l == INFO {
		return "info"
	} else if l == SUCCESS {
		return "success"
	} else if l == WARN {
		return "warn"
	} else if l == ERROR {
		return "error"
	}
	return ""
}

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
	Timestamp    time.Time
	LogKey       string
	LogValues    map[string]string
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
