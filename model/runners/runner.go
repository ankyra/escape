package runners

import (
	"fmt"
	"os"

	core "github.com/ankyra/escape-core"
	"github.com/ankyra/escape-core/state"
	"github.com/ankyra/escape/model/paths"
)

type Runner interface {
	Run(*RunnerContext) error
}

type runner struct {
	run func(ctx *RunnerContext) error
}

func (r *runner) Run(ctx *RunnerContext) error {
	return r.run(ctx)
}

func NewRunner(r func(ctx *RunnerContext) error) Runner {
	return &runner{
		run: r,
	}
}

func NewCompoundRunner(runners ...Runner) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		for _, r := range runners {
			if err := r.Run(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func NewStatusCodeRunner(stage string, status state.StatusCode) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		st := state.NewStatus(status)
		return ctx.GetDeploymentState().UpdateStatus(stage, st)
	})
}

func NewDependencyRunner(logKey, stage string, depRunner func() Runner, errorCode state.StatusCode) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		parentInputs, err := NewEnvironmentBuilder().GetPreDependencyInputs(ctx, stage)
		if err != nil {
			return ReportFailure(ctx, stage, err, errorCode)
		}
		metadata := ctx.GetReleaseMetadata()
		for _, depend := range metadata.Depends {
			if err := runDependency(ctx, depend, logKey, stage, depRunner(), parentInputs); err != nil {
				return ReportFailure(ctx, stage, err, errorCode)
			}
		}
		return nil
	})
}

func NewProviderActivationRunner(stage string) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		metadata := ctx.GetReleaseMetadata()
		if err := ctx.deploymentState.ConfigureProviders(metadata, stage, nil); err != nil {
			return err
		}
		for _, consume := range metadata.Consumes {
			if consume.InScope(stage) {
				return runProvider(stage, "activate", ctx, consume)
			}
		}
		return nil
	})
}

func NewProviderDeactivationRunner(stage string) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		metadata := ctx.GetReleaseMetadata()
		if err := ctx.deploymentState.ConfigureProviders(metadata, stage, nil); err != nil {
			return err
		}
		for _, consume := range metadata.Consumes {
			if consume.InScope(stage) {
				return runProvider(stage, "deactivate", ctx, consume)
			}
		}
		return nil
	})
}

func runProvider(stage, action string, ctx *RunnerContext, consume *core.ConsumerConfig) error {
	ctx.Logger().Log("provider."+action, map[string]string{
		"variable": consume.VariableName,
		"consumes": consume.Name,
	})
	deploymentName := ctx.deploymentState.GetStageOrCreateNew(stage).Providers[consume.VariableName]
	depl, found := ctx.environmentState.Deployments[deploymentName]
	if !found {
		return fmt.Errorf("Deployment %s which is the configured provider for '%s' ($%s) could not be found",
			deploymentName, consume.Name, consume.VariableName)
	}
	releaseId := depl.GetReleaseId("deploy")
	ctx.Logger().PushSection("Provider " + releaseId + ", implements " + consume.Name)
	ctx.Logger().PushRelease(releaseId)

	// TODO: move into provider directory
	// TODO: prepare a ScriptRunner
	// TODO: call activate/deactivate script
	// TODO: skip if script is not set

	ctx.Logger().PopRelease()
	ctx.Logger().PopSection()
	return nil
}

func runDependency(ctx *RunnerContext, depCfg *core.DependencyConfig, logKey, stage string, runner Runner, parentInputs map[string]interface{}) error {
	dependency := depCfg.ReleaseId
	ctx.Logger().PushSection("Dependency " + dependency)
	ctx.Logger().Log(logKey+"."+logKey+"_dependency", map[string]string{
		"dependency": dependency,
	})
	ctx.Logger().PushRelease(dependency)
	mapping := depCfg.GetMapping(stage)
	if mapping == nil {
		return fmt.Errorf("Invalid stage '%s'", stage)
	}
	inputs, err := NewEnvironmentBuilder().GetInputsForDependency(ctx, stage, mapping, parentInputs)
	if err != nil {
		return err
	}
	dep, err := core.NewDependencyFromString(dependency)
	if err != nil {
		return err
	}
	location := ctx.GetPath().UnpackedDepDirectory(dep)
	metadata, err := newMetadataFromReleaseDir(location)
	if err != nil {
		return err
	}
	depCtx, err := ctx.NewContextForDependency(depCfg.DeploymentName, metadata, depCfg.Consumes)
	if err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	deplState := depCtx.GetDeploymentState()
	if err := deplState.UpdateUserInputs("deploy", inputs); err != nil {
		return err
	}
	if err := os.Chdir(location); err != nil {
		return err
	}
	if err := runner.Run(depCtx); err != nil {
		return err
	}
	if err := os.Chdir(currentDir); err != nil {
		return err
	}
	ctx.Logger().Log(logKey+"."+logKey+"_dependency_finished", nil)
	ctx.Logger().PopRelease()
	ctx.Logger().PopSection()
	return nil
}

func newMetadataFromReleaseDir(releaseDir string) (*core.ReleaseMetadata, error) {
	path := paths.NewPathWithBaseDir(releaseDir).ReleaseJson()
	return core.NewReleaseMetadataFromFile(path)
}
