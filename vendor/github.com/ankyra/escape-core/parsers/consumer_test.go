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

package parsers

import (
	"fmt"

	. "gopkg.in/check.v1"
)

type consumerSuite struct{}

var _ = Suite(&consumerSuite{})

func (s *consumerSuite) Test_ParseConsumer_Happy_Path(c *C) {
	cases := []string{
		"provider",
		"provider  ",
		"  provider",
		"   provider  ",
	}
	renameCases := []string{
		"provider as p1",
		"provider  as p1",
		"provider as  p1",
		"provider as p1 ",
		"provider  as  p1",
		"  provider as p1",
		"  provider  as p1",
		"  provider  as  p1",
		"  provider  as  p1  ",
	}
	for _, test := range cases {
		consumer, err := ParseConsumer(test)
		c.Assert(err, IsNil)
		c.Assert(consumer.Interface, Equals, "provider")
		c.Assert(consumer.VariableName, Equals, "")
	}
	for _, test := range renameCases {
		consumer, err := ParseConsumer(test)
		c.Assert(err, IsNil)
		c.Assert(consumer.Interface, Equals, "provider")
		c.Assert(consumer.VariableName, Equals, "p1")
	}
}

func (s *consumerSuite) Test_ParseConsumer_fails_when_wrong_number_of_parts(c *C) {
	cases := []string{
		"",
		"provider as",
		"provider as p1 many",
		"provider as p1 asdopjkasdpok asdpoka sdpokas dpokas dpok",
	}
	for _, test := range cases {
		_, err := ParseConsumer(test)
		c.Assert(err, DeepEquals, fmt.Errorf("Malformed consumer string '%s'", test))
	}
}

func (s *consumerSuite) Test_ParseConsumer_fails_when_missing_as(c *C) {
	cases := []string{
		"provider wut p1",
	}
	for _, test := range cases {
		_, err := ParseConsumer(test)
		c.Assert(err, DeepEquals, fmt.Errorf("Unexpected 'wut' expecting 'as' in '%s'", test))
	}
}

func (s *consumerSuite) Test_ParseConsumer_fails_with_invalid_variable_name(c *C) {
	cases := []string{
		"provider as $23",
		"provider  as $23",
		"provider as  $23",
		"provider as $23 ",
		"provider  as  $23",
		"  provider as $23",
		"  provider  as $23",
		"  provider  as  $23",
		"  provider  as  $23  ",
	}
	for _, test := range cases {
		_, err := ParseConsumer(test)
		c.Assert(err, DeepEquals, fmt.Errorf("Malformed consumer string '%s': Invalid variable format '$23'", test))
	}
}
