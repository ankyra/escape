package godog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/ankyra/escape-client/model/escape_plan"
	"github.com/ankyra/escape-client/model/state"
	eutil "github.com/ankyra/escape-client/util"
	state_types "github.com/ankyra/escape-core/state"
	"github.com/ankyra/escape-core/util"
	"gopkg.in/yaml.v2"
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
		os.RemoveAll("deps/")
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

func inputVariableWithDefaultAndItems(variableId, defaultValue, items string) error {
	plan := escape_plan.NewEscapePlan()
	err := plan.LoadConfig("escape.yml")
	if err != nil {
		return nil
	}
	plan.Inputs = append(plan.Inputs, map[string]interface{}{
		"id":      variableId,
		"default": defaultValue,
		"items":   items,
	})
	return ioutil.WriteFile("escape.yml", plan.ToMinifiedYaml(), 0644)
}

func outputVariableWithDefault(variableId, defaultValue string) error {
	return outputVariable("string", variableId, defaultValue)
}

func outputVariable(typ, variableId, defaultValue string) error {
	plan := escape_plan.NewEscapePlan()
	err := plan.LoadConfig("escape.yml")
	if err != nil {
		return nil
	}
	plan.Outputs = append(plan.Outputs, map[string]interface{}{
		"id":      variableId,
		"default": defaultValue,
		"type":    typ,
	})
	return ioutil.WriteFile("escape.yml", plan.ToMinifiedYaml(), 0644)
}

func outputListVariableWithDefault(variableId, defaultValue string) error {
	return outputVariable("list[string]", variableId, defaultValue)
}

func templateContaining(filename, content string) error {
	return templateContainingWithScope(filename, content, "")
}

func templateContainingWithScope(filename, content, scope string) error {
	scopes := []string{"build", "deploy"}
	if scope != "" {
		scopes = []string{scope}
	}
	plan := escape_plan.NewEscapePlan()
	if err := plan.LoadConfig("escape.yml"); err != nil {
		return nil
	}
	plan.Templates = append(plan.Templates, map[string]interface{}{
		"file":   filename,
		"scopes": scopes,
	})
	if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
		return err
	}
	return ioutil.WriteFile("escape.yml", plan.ToMinifiedYaml(), 0644)
}

func iShouldNotHaveAFile(arg1 string) error {
	_, err := os.Stat(arg1)
	if err == nil {
		return fmt.Errorf("File '%s' exists", arg1)
	}
	return nil
}

func iShouldHaveAFileWithContents(filename, expectedContent string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if string(bytes) != expectedContent {
		return fmt.Errorf("Expecting '%s' got '%s'", expectedContent, string(bytes))
	}
	return nil
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
		return fmt.Errorf("'%s' not found in calculated inputs", key)
	}
	if v != value {
		return fmt.Errorf("Expecting '%s', got '%s'", value, v)
	}
	return nil
}

func itsCalculatedOutputIsSetTo(key, value string) error {
	inputs := CapturedDeployment.GetCalculatedOutputs(CapturedStage)
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

func iDeploy(arg1 string) error {
	rec := util.NewProcessRecorder()
	cmd := []string{"escape", "deploy", arg1}
	stdout, err := rec.Record(cmd, nil, eutil.NewLoggerDummy())
	CapturedStdout = stdout
	if err != nil {
		fmt.Println(stdout)
	}
	CapturedStage = "deploy"
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

func itProvides(arg1 string) error {
	plan := escape_plan.NewEscapePlan()
	err := plan.LoadConfig("escape.yml")
	if err != nil {
		return nil
	}
	plan.Provides = append(plan.Provides, arg1)
	return ioutil.WriteFile("escape.yml", plan.ToMinifiedYaml(), 0644)
}

func itConsumes(arg1 string) error {
	plan := escape_plan.NewEscapePlan()
	err := plan.LoadConfig("escape.yml")
	if err != nil {
		return nil
	}
	plan.Consumes = append(plan.Consumes, arg1)
	return ioutil.WriteFile("escape.yml", plan.ToMinifiedYaml(), 0644)
}

func isTheProviderFor(deploymentName, providerName string) error {
	d := CapturedDeployment.GetDeploymentOrMakeNew("build", deploymentName)
	prov, found := d.GetProviders("build")[providerName]
	if !found {
		return fmt.Errorf("'%s' provider not found", providerName)
	}
	if prov != deploymentName {
		return fmt.Errorf("'%s' provider is '%s' not expected '%s'", providerName, prov, deploymentName)
	}
	return nil
}

func versionIsPresentInItsDeploymentState(deploymentName, version string) error {
	d := CapturedDeployment.GetDeploymentOrMakeNew("build", deploymentName)
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

func versionIsPresentInTheDeployState(deploymentName, version string) error {
	env, err := state.NewLocalStateProvider("escape_state.json").Load("prj", "dev")
	if err != nil {
		return err
	}
	d := env.GetOrCreateDeploymentState(deploymentName)
	if d.GetVersion("deploy") != version {
		return fmt.Errorf("Expecting '%s', got '%s'", version, d.GetVersion("build"))
	}
	CapturedDeployment = d
	CapturedStage = "deploy"
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
	s.Step(`^I build the application again$`, iBuildTheApplication)
	s.Step(`^I deploy "([^"]*)"$`, iDeploy)
	s.Step(`^"([^"]*)" version "([^"]*)" is present in the build state$`, versionIsPresentInTheBuildState)
	s.Step(`^"([^"]*)" version "([^"]*)" is present in the deploy state$`, versionIsPresentInTheDeployState)
	s.Step(`^input variable "([^"]*)" with default "([^"]*)"$`, inputVariableWithDefault)
	s.Step(`^input variable "([^"]*)" with default "([^"]*)" and items "([^"]*)"$`, inputVariableWithDefaultAndItems)
	s.Step(`^its calculated input "([^"]*)" is set to "([^"]*)"$`, itsCalculatedInputIsSetTo)
	s.Step(`^its calculated output "([^"]*)" is set to "([^"]*)"$`, itsCalculatedOutputIsSetTo)
	s.Step(`^I release the application$`, iReleaseTheApplication)
	s.Step(`^it has "([^"]*)" as a dependency$`, itHasAsADependency)
	s.Step(`^it has "([^"]*)" set to "([^"]*)"$`, itHasSetTo)
	s.Step(`^"([^"]*)" version "([^"]*)" is present in its deployment state$`, versionIsPresentInItsDeploymentState)
	s.Step(`^it provides "([^"]*)"$`, itProvides)
	s.Step(`^it consumes "([^"]*)"$`, itConsumes)
	s.Step(`^"([^"]*)" is the provider for "([^"]*)"$`, isTheProviderFor)
	s.Step(`^output variable "([^"]*)" with default "([^"]*)"$`, outputVariableWithDefault)
	s.Step(`^output variable "([^"]*)" with default '([^']*)'$`, outputVariableWithDefault)
	s.Step(`^list output variable "([^"]*)" with default '([^']*)'$`, outputListVariableWithDefault)
	s.Step(`^template "([^"]*)" containing "([^"]*)" with "([^"]*)" scope$`, templateContainingWithScope)
	s.Step(`^I should not have a file "([^"]*)"$`, iShouldNotHaveAFile)
	s.Step(`^template "([^"]*)" containing "([^"]*)"$`, templateContaining)
	s.Step(`^I should have a file "([^"]*)" with contents "([^"]*)"$`, iShouldHaveAFileWithContents)

	s.BeforeScenario(func(interface{}) {
		StartRegistry()
	})
	s.AfterScenario(func(interface{}, error) {
		StopRegistry()
	})
}
