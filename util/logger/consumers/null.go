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

import "github.com/ankyra/escape/util/logger/api"

type nullLogConsumer struct{}

func NewNullLogConsumer() *nullLogConsumer {
	return &nullLogConsumer{}
}

func (t *nullLogConsumer) Consume(entry *api.LogEntry) error {
	return nil
}

func (t *nullLogConsumer) Close() {
}
