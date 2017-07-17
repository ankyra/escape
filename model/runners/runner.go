package runners

import (
	"github.com/ankyra/escape-client/model/paths"
	core "github.com/ankyra/escape-core"
	"os"
)

type Runner interface {
	Run(RunnerContext) error
}

type runner struct {
	run func(ctx RunnerContext) error
}

func (r *runner) Run(ctx RunnerContext) error {
	return r.run(ctx)
}

func NewRunner(r func(ctx RunnerContext) error) Runner {
	return &runner{
		run: r,
	}
}

func NewCompoundRunner(runners ...Runner) Runner {
	return NewRunner(func(ctx RunnerContext) error {
		for _, r := range runners {
			if err := r.Run(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func NewDependencyRunner(logKey string, depRunner func() Runner) Runner {
	return NewRunner(func(ctx RunnerContext) error {
        inputs, err := NewEnvironmentBuilder().GetPreDependencyInputs(ctx, stage)
        if err != nil {
            return err
        }
		metadata := ctx.GetReleaseMetadata()
		for _, depend := range metadata.Depends {
            // TODO pass in stage
            // TODO use inputs for GetInputsForDependency
			if err := runDependency(ctx, depend, logKey, depRunner(), inputs); err != nil {
				return err
			}
		}
		return nil
	})
}

func runDependency(ctx RunnerContext, depCfg *core.DependencyConfig, logKey string, runner Runner) error {
	dependency := depCfg.ReleaseId
	ctx.Logger().PushSection("Dependency " + dependency)
	ctx.Logger().Log(logKey+"."+logKey+"_dependency", map[string]string{
		"dependency": dependency,
	})
	ctx.Logger().PushRelease(dependency)
    inputs, err := NewEnvironmentBuilder().GetInputsForDependency(ctx, depCfg)
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
	depCtx := ctx.NewContextForDependency(metadata)
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
