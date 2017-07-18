package main

import (
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/ankyra/escape-client/model/escape_plan"
	"github.com/ankyra/escape-client/model/state"
	state_types "github.com/ankyra/escape-client/model/state/types"
	eutil "github.com/ankyra/escape-client/util"
	"github.com/ankyra/escape-core/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

var CapturedStdout string
var CapturedDeployment *state_types.DeploymentState
var CapturedStage string
var ServerProcess *exec.Cmd

func StartRegistry() {
	go func() {
		os.RemoveAll("test.db")
		os.RemoveAll("escape_state.json")
		os.RemoveAll("escape.yml")
		os.RemoveAll("releases/")
		os.Mkdir("releases/", 0755)
		env := []string{
			"DATABASE=sqlite",
			"DATABASE_SETTINGS_PATH=test.db",
			"STORAGE_BACKEND=local",
			"STORAGE_SETTINGS_PATH=releases/",
		}
		ServerProcess = exec.Command("escape-registry")
		ServerProcess.Env = env
		if err := ServerProcess.Start(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second / 2)
}

func StopRegistry() {
	ServerProcess.Process.Kill()
}

func aNewEscapePlanCalled(name string) error {
	rec := util.NewProcessRecorder()
	cmd := []string{"escape", "plan", "init", "-f", "-n", name}
	stdout, err := rec.Record(cmd, nil, eutil.NewLoggerDummy())
	CapturedStdout = stdout
	return err
}

func iCompileThePlan() error {
	rec := util.NewProcessRecorder()
	cmd := []string{"escape", "plan", "compile"}
	stdout, err := rec.Record(cmd, nil, eutil.NewLoggerDummy())
	CapturedStdout = stdout
	return err
}

func inputVariableWithDefault(variableId, defaultValue string) error {
	plan := escape_plan.NewEscapePlan()
	err := plan.LoadConfig("escape.yml")
	if err != nil {
		return nil
	}
	plan.Inputs = append(plan.Inputs, map[string]interface{}{
		"id":      variableId,
		"default": defaultValue,
	})
	return ioutil.WriteFile("escape.yml", plan.ToMinifiedYaml(), 0644)
}

func itHasSetTo(arg1, arg2 string) error {
	plan := map[string]interface{}{}
	bytes, err := ioutil.ReadFile("escape.yml")
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, &plan); err != nil {
		return err
	}
	plan[arg1] = arg2
	bytes, err = yaml.Marshal(plan)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("escape.yml", bytes, 0644)
}

func itsCalculatedInputIsSetTo(key, value string) error {
	inputs := CapturedDeployment.GetCalculatedInputs(CapturedStage)
	v, found := inputs[key]
	if !found {
		return fmt.Errorf("'%s' not found in calculated outputs", key)
	}
	if v != value {
		return fmt.Errorf("Expecting '%s', got '%s'", value, v)
	}
	return nil
}

func iBuildTheApplication() error {
	rec := util.NewProcessRecorder()
	cmd := []string{"escape", "build"}
	stdout, err := rec.Record(cmd, nil, eutil.NewLoggerDummy())
	CapturedStdout = stdout
	if err != nil {
		fmt.Println(stdout)
	}
	return err
}

func iReleaseTheApplication() error {
	rec := util.NewProcessRecorder()
	cmd := []string{"escape", "release", "-f"}
	stdout, err := rec.Record(cmd, nil, eutil.NewLoggerDummy())
	CapturedStdout = stdout
	if err != nil {
		fmt.Println(stdout)
	}
	return err
}

func itHasAsADependency(dependency string) error {
	plan := escape_plan.NewEscapePlan()
	err := plan.LoadConfig("escape.yml")
	if err != nil {
		return nil
	}
	plan.Depends = append(plan.Depends, dependency)
	return ioutil.WriteFile("escape.yml", plan.ToMinifiedYaml(), 0644)
}

func versionIsPresentInItsDeploymentState(deploymentName, version string) error {
	d := CapturedDeployment.GetDeployment("build", deploymentName)
	if d.GetVersion("deploy") != version {
		return fmt.Errorf("Expecting '%s', got '%s'", version, d.GetVersion("deploy"))
	}
	CapturedDeployment = d
	CapturedStage = "deploy"
	return nil
}

func versionIsPresentInTheBuildState(deploymentName, version string) error {
	env, err := state.NewLocalStateProvider("escape_state.json").Load("prj", "dev")
	if err != nil {
		return err
	}
	d := env.GetOrCreateDeploymentState(deploymentName)
	if d.GetVersion("build") != version {
		return fmt.Errorf("Expecting '%s', got '%s'", version, d.GetVersion("build"))
	}
	CapturedDeployment = d
	CapturedStage = "build"
	return nil
}

func iShouldHaveValidReleaseMetadata() error {
	result := map[string]interface{}{}
	err := json.Unmarshal([]byte(CapturedStdout), &result)
	if err != nil {
		return err
	}
	return nil
}

func theMetadataShouldHaveItsSetTo(key, value string) error {
	result := map[string]interface{}{}
	err := json.Unmarshal([]byte(CapturedStdout), &result)
	if err != nil {
		return err
	}
	v, found := result[key]
	if !found {
		return fmt.Errorf("'%s' not found in metadata", key)
	}
	if v != value {
		return fmt.Errorf("Expecting '%s', got '%s'", value, v)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^a new Escape plan called "([^"]*)"$`, aNewEscapePlanCalled)
	s.Step(`^I compile the plan$`, iCompileThePlan)
	s.Step(`^I should have valid release metadata$`, iShouldHaveValidReleaseMetadata)
	s.Step(`^the metadata should have its "([^"]*)" set to "([^"]*)"$`, theMetadataShouldHaveItsSetTo)
	s.Step(`^I build the application$`, iBuildTheApplication)
	s.Step(`^"([^"]*)" version "([^"]*)" is present in the build state$`, versionIsPresentInTheBuildState)
	s.Step(`^input variable "([^"]*)" with default "([^"]*)"$`, inputVariableWithDefault)
	s.Step(`^its calculated input "([^"]*)" is set to "([^"]*)"$`, itsCalculatedInputIsSetTo)
	s.Step(`^I release the application$`, iReleaseTheApplication)
	s.Step(`^it has "([^"]*)" as a dependency$`, itHasAsADependency)
	s.Step(`^it has "([^"]*)" set to "([^"]*)"$`, itHasSetTo)
	s.Step(`^"([^"]*)" version "([^"]*)" is present in its deployment state$`, versionIsPresentInItsDeploymentState)
	s.BeforeScenario(func(interface{}) {
		StartRegistry()
	})
	s.AfterScenario(func(interface{}, error) {
		StopRegistry()
	})
}
