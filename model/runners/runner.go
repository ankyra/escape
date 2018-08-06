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

func NewDependencyRunner(logKey, parentStage string, depRunner func() Runner, errorCode state.StatusCode) Runner {
	return NewRunner(func(ctx *RunnerContext) error {
		parentInputs, err := NewEnvironmentBuilder().GetPreDependencyInputs(ctx, parentStage)
		if err != nil {
			return ReportFailure(ctx, parentStage, err, errorCode)
		}
		metadata := ctx.GetReleaseMetadata()
		for _, depend := range metadata.Depends {
			if err := runDependency(ctx, depend, logKey, parentStage, depRunner(), parentInputs); err != nil {
				return ReportFailure(ctx, parentStage, err, errorCode)
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
	deploymentPath := ctx.deploymentState.GetStageOrCreateNew(stage).Providers[consume.VariableName]

	depl, err := ctx.environmentState.ResolveDeploymentPath(ctx.deploymentState.GetRootDeploymentStage(), deploymentPath)
	if err != nil {
		return fmt.Errorf("Failed to load configured provider for '%s' ($%s), root stage %s, path %s: %s",
			consume.Name, consume.VariableName, ctx.deploymentState.GetRootDeploymentStage(), deploymentPath, err.Error())
	}
	releaseId := depl.GetReleaseId("deploy")
	ctx.Logger().PushSection("Provider " + releaseId + "($" + consume.Name + ")")
	ctx.Logger().PushRelease(releaseId)

	depCfg := core.NewDependencyConfig(releaseId)
	metadata, err := ctx.context.GetDependencyMetadata(depCfg)
	if err != nil {
		return err
	}

	execStage := metadata.GetExecStage(action + "_provider")
	if execStage == nil {
		return nil
	}

	dep, err := core.NewDependencyFromString(releaseId)
	if err != nil {
		return err
	}
	location := ctx.GetPath().UnpackedDepDirectory(dep)

	ctx.Logger().Log("provider."+action, map[string]string{
		"variable": consume.VariableName,
		"consumes": consume.Name,
	})

	newCtx, err := ctx.NewContextForProvider(depl, metadata)
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(location); err != nil {
		return err
	}
	runner := NewScriptRunner("deploy", action+"_provider", state.OK, state.Failure)
	if err := runner.Run(newCtx); err != nil {
		return err
	}
	if err := os.Chdir(currentDir); err != nil {
		return err
	}

	ctx.Logger().Log("provider."+action+".finished", map[string]string{
		"variable": consume.VariableName,
		"consumes": consume.Name,
	})

	ctx.Logger().PopRelease()
	ctx.Logger().PopSection()
	return nil
}

func runDependency(ctx *RunnerContext, depCfg *core.DependencyConfig, logKey, parentStage string, runner Runner, parentInputs map[string]interface{}) error {
	dependency := depCfg.ReleaseId
	ctx.Logger().PushSection("Dependency " + dependency)
	ctx.Logger().Log(logKey+"."+logKey+"_dependency", map[string]string{
		"dependency": dependency,
	})
	ctx.Logger().PushRelease(dependency)
	mapping := depCfg.GetMapping(parentStage)
	if mapping == nil {
		return fmt.Errorf("Invalid stage '%s'", parentStage)
	}
	inputs, err := NewEnvironmentBuilder().GetInputsForDependency(ctx, parentStage, mapping, parentInputs)
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
	depCtx, err := ctx.NewContextForDependency(parentStage, depCfg.DeploymentName, metadata, depCfg.Consumes)
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
