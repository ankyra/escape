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

package consumers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/ankyra/escape/util/logger/api"
)

type jsonLogConsumer struct {
	PreviousSectionStack []string
	Silent               bool
}

func NewJSONLogConsumer() *jsonLogConsumer {
	return &jsonLogConsumer{
		Silent: false,
	}
}

type JSONMessage struct {
	Timestamp         time.Time         `json:"timestamp"`
	Message           string            `json:"message"`
	Level             string            `json:"level"`
	LogKey            string            `json:"log_key"`
	LogValues         map[string]string `json:"log_values"`
	LogSectionChanged bool              `json:"log_section_changed"`
	LogSections       []string          `json:"log_sections"`
}

func (t *jsonLogConsumer) Consume(entry *api.LogEntry) (string, error) {
	sectionChanged := !reflect.DeepEqual(t.PreviousSectionStack, entry.SectionStack)
	msg := JSONMessage{
		Timestamp:         entry.Timestamp,
		Message:           entry.Message,
		Level:             entry.LogLevel.String(),
		LogKey:            entry.LogKey,
		LogValues:         entry.LogValues,
		LogSectionChanged: sectionChanged,
		LogSections:       entry.SectionStack,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	str := string(bytes)
	if !t.Silent {
		fmt.Println(str)
	}
	return str, nil
}
