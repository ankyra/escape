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
	"time"

	"github.com/ankyra/escape/util/logger/api"
)

type jsonLogConsumer struct{}

func NewJSONLogConsumer() *jsonLogConsumer {
	return &jsonLogConsumer{}
}

type JSONMessage struct {
	Timestamp time.Time `json:"timestampe"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
}

func (t *jsonLogConsumer) Consume(entry *api.LogEntry) error {
	return nil
}
