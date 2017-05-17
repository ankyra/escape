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

package util

import (
	"bytes"
	"text/template"
)

type LogLevel int

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
}

type logger struct {
	sections  []string
	releases  []string
	consumers []LogConsumer
}

func NewLogger(consumers []LogConsumer) Logger {
	return &logger{
		sections:  []string{},
		consumers: consumers,
	}
}

func (l *logger) Log(key string, values map[string]string) {
	msg, ok := logMessages[key]
	if !ok {
		panic("Unknown log key: " + key)
	}
	tpl, err := template.New("tpl").Parse(msg["msg"])
	if err != nil {
		panic(err)
	}
	if values == nil {
		values = map[string]string{}
	}
	if len(l.releases) > 0 {
		values["release"] = l.releases[len(l.releases)-1]
	}
	writer := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(writer, values); err != nil {
		panic(err)
	}
	level := INFO
	if msg["level"] == "debug" {
		level = DEBUG
	} else if msg["level"] == "warn" {
		level = WARN
	} else if msg["level"] == "success" {
		level = SUCCESS
	} else if msg["level"] == "error" {
		level = ERROR
	}
	collapse, ok := msg["collapse"]
	if !ok {
		collapse = "true"
	}
	entry := &LogEntry{
		Message:      writer.String(),
		SectionStack: l.sections,
		LogLevel:     LogLevel(level),
		Collapse:     collapse == "true",
	}
	for _, c := range l.consumers {
		if err := c.Consume(entry); err != nil {
			panic(err)
		}
	}
}

func (l *logger) PushSection(s string) {
	l.sections = append(l.sections, s)
}

func (l *logger) PopSection() {
	l.sections = l.sections[:len(l.sections)-1]
}

func (l *logger) PushRelease(s string) {
	l.releases = append(l.releases, s)
}

func (l *logger) PopRelease() {
	l.releases = l.releases[:len(l.releases)-1]
}
