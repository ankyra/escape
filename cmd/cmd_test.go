package cmd

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type suite struct{}

var _ = Suite(&suite{})

func (s *suite) SetUpTest(c *C) {
	os.Remove("escape.yml")
	readLocalErrands = false
	escapePlanLocation = "escape.yml"
	state = "escape_state.json"
	environment = "dev"
	deployment = ""
}
