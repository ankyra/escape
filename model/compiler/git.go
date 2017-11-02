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

package compiler

import (
	"os"

	"github.com/ankyra/escape/util"
)

func compileGit(ctx *CompilerContext) error {
	rec := util.NewProcessRecorder()
	result, err := rec.Record([]string{"git", "rev-parse", "HEAD"}, os.Environ(), util.NewLoggerDummy())
	if err != nil {
		return nil
	}
	ctx.Metadata.Revision = result

	result, err = rec.Record([]string{"git", "rev-parse", "--abbrev-ref", "HEAD"}, os.Environ(), util.NewLoggerDummy())
	if err != nil {
		return nil
	}
	ctx.Metadata.Branch = result

	result, err = rec.Record([]string{"git", "config", "--get", "remote.origin.url"}, os.Environ(), util.NewLoggerDummy())
	if err != nil {
		return nil
	}
	ctx.Metadata.Repository = result

	result, err = rec.Record([]string{"git", "show", "-s", "--format=%B", "HEAD"}, os.Environ(), util.NewLoggerDummy())
	if err != nil {
		return nil
	}
	ctx.Metadata.RevisionMessage = result

	result, err = rec.Record([]string{"git", "show", "-s", "--format=%an <%ae>", "HEAD"}, os.Environ(), util.NewLoggerDummy())
	if err != nil {
		return nil
	}
	ctx.Metadata.RevisionAuthor = result
	return nil
}
