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

package core

import (
	. "gopkg.in/check.v1"
)

func (s *metadataSuite) Test_NewConsumerConfig(c *C) {
	consumer := NewConsumerConfig("test")
	c.Assert(consumer.Name, Equals, "test")
	c.Assert(consumer.Scopes, DeepEquals, []string{"build", "deploy"})
	c.Assert(consumer.VariableName, Equals, "test")
}

func (s *metadataSuite) Test_NewConsumerConfigFromMap(c *C) {
	consumer, err := NewConsumerConfigFromMap(map[interface{}]interface{}{
		"name":   "test",
		"scopes": []interface{}{"build", "deploy"},
	})
	c.Assert(err, IsNil)
	c.Assert(consumer.Name, Equals, "test")
	c.Assert(consumer.Scopes, DeepEquals, []string{"build", "deploy"})
	c.Assert(consumer.VariableName, Equals, "test")
}

func (s *metadataSuite) Test_NewConsumerConfigFromMap_renamed_var(c *C) {
	consumer, err := NewConsumerConfigFromMap(map[interface{}]interface{}{
		"name":   "test as t",
		"scopes": []interface{}{"build", "deploy"},
	})
	c.Assert(err, IsNil)
	c.Assert(consumer.Name, Equals, "test")
	c.Assert(consumer.Scopes, DeepEquals, []string{"build", "deploy"})
	c.Assert(consumer.VariableName, Equals, "t")
}

func (s *metadataSuite) Test_NewConsumerConfigFromMap_No_Scopes_360_blaze_it(c *C) {
	consumer, err := NewConsumerConfigFromMap(map[interface{}]interface{}{
		"name": "test",
	})
	c.Assert(err, IsNil)
	c.Assert(consumer.Name, Equals, "test")
	c.Assert(consumer.Scopes, DeepEquals, []string{"build", "deploy"})
	c.Assert(consumer.VariableName, Equals, "test")
}

func (s *metadataSuite) Test_NewConsumerConfigFromMap_limited_scope(c *C) {
	consumer, err := NewConsumerConfigFromMap(map[interface{}]interface{}{
		"name":   "test",
		"scopes": []interface{}{"deploy"},
	})
	c.Assert(err, IsNil)
	c.Assert(consumer.Name, Equals, "test")
	c.Assert(consumer.Scopes, DeepEquals, []string{"deploy"})
	c.Assert(consumer.VariableName, Equals, "test")
}

func (s *metadataSuite) Test_NewConsumerConfigFromInterface_String(c *C) {
	consumer, err := NewConsumerConfigFromInterface("test")
	c.Assert(err, IsNil)
	c.Assert(consumer.Name, Equals, "test")
	c.Assert(consumer.Scopes, DeepEquals, []string{"build", "deploy"})
	c.Assert(consumer.VariableName, Equals, "test")
}

func (s *metadataSuite) Test_NewConsumerConfigFromInterface_Renamed_String(c *C) {
	consumer, err := NewConsumerConfigFromInterface("test as t")
	c.Assert(err, IsNil)
	c.Assert(consumer.Name, Equals, "test")
	c.Assert(consumer.Scopes, DeepEquals, []string{"build", "deploy"})
	c.Assert(consumer.VariableName, Equals, "t")
}

func (s *metadataSuite) Test_NewConsumerConfigFromInterface_Renamed_String_fails_if_invalid_Format(c *C) {
	cases := []string{
		"test as $23",
		"test this p1",
		"test as   $23",
		"",
	}
	for _, test := range cases {
		_, err := NewConsumerConfigFromInterface(test)
		c.Assert(err, Not(IsNil))
	}
}

func (s *metadataSuite) Test_NewConsumerConfigFromInterface_Map(c *C) {
	consumer, err := NewConsumerConfigFromInterface(map[interface{}]interface{}{"name": "test"})
	c.Assert(err, IsNil)
	c.Assert(consumer.Name, Equals, "test")
	c.Assert(consumer.Scopes, DeepEquals, []string{"build", "deploy"})
	c.Assert(consumer.VariableName, Equals, "test")
}

func (s *metadataSuite) Test_NewConsumerConfigFromInterface_fails_on_wrong_type(c *C) {
	_, err := NewConsumerConfigFromInterface(12)
	c.Assert(err.Error(), Equals, "Expecting dict or string type")
}

func (s *metadataSuite) Test_ConsumerConfig_Validate_sets_scopes_if_nil(c *C) {
	consumer := NewConsumerConfig("test")
	consumer.Scopes = nil
	c.Assert(consumer.ValidateAndFix(), IsNil)
	c.Assert(consumer.Scopes, DeepEquals, []string{"build", "deploy"})
}

func (s *metadataSuite) Test_ConsumerConfig_Validate_sets_scopes_if_empty(c *C) {
	consumer := NewConsumerConfig("test")
	consumer.Scopes = []string{}
	c.Assert(consumer.ValidateAndFix(), IsNil)
	c.Assert(consumer.Scopes, DeepEquals, []string{"build", "deploy"})
}
