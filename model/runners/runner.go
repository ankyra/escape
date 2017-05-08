package runners

import (
	"github.com/ankyra/escape-client/model"
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/release"
	"os"
)

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

func NewDependencyRunner(stage string, depRunner func() Runner) Runner {
	return NewRunner(func(ctx RunnerContext) error {
		metadata := ctx.GetReleaseMetadata()
		for _, depend := range metadata.GetDependencies() {
			if err := runDependency(ctx, depend, stage, depRunner()); err != nil {
				return err
			}
		}
		return nil
	})
}

func runDependency(ctx RunnerContext, dependency, stage string, runner Runner) error {
	ctx.Logger().PushSection("Dependency " + dependency)
	ctx.Logger().Log(stage+"."+stage+"_dependency", map[string]string{
		"dependency": dependency,
	})
	ctx.Logger().PushRelease(dependency)
	dep, err := release.NewDependencyFromString(dependency)
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
	if err := os.Chdir(location); err != nil {
		return err
	}
	if err := runner.Run(depCtx); err != nil {
		return err
	}
	if err := os.Chdir(currentDir); err != nil {
		return err
	}
	ctx.Logger().Log(stage+"."+stage+"_dependency_finished", nil)
	ctx.Logger().PopRelease()
	ctx.Logger().PopSection()
	return nil
}

func newMetadataFromReleaseDir(releaseDir string) (ReleaseMetadata, error) {
	path := model.NewPathWithBaseDir(releaseDir).ReleaseJson()
	return release.NewReleaseMetadataFromFile(path)
}
