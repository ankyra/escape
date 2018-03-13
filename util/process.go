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

package util

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ankyra/escape/util/logger/api"
)

type ProcessRecorder interface {
	SetWorkingDirectory(string)
	Record(cmd []string, env []string, log api.Logger) (string, error)
	Run(cmd []string, env []string, log api.Logger) error
}

type processRecorder struct {
	WorkingDirectory string
}

func NewProcessRecorder() ProcessRecorder {
	return &processRecorder{}
}

func (p *processRecorder) SetWorkingDirectory(cwd string) {
	p.WorkingDirectory = cwd
}

func getExtraPathDir() string {
	currentUser, _ := user.Current()
	return filepath.Join(GetAppConfigDir(runtime.GOOS, currentUser.HomeDir), ".bin")
}

func pipeReader(pipe io.ReadCloser, channel chan string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		channel <- line
	}
	close(channel)
}

func (p *processRecorder) Record(cmd []string, env []string, log api.Logger) (string, error) {
	extraPath := getExtraPathDir()
	MkdirRecursively(extraPath)
	newEnv := []string{}
	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			path := e[5:]
			newEnv = append(newEnv, "PATH="+extraPath+":"+path)
		} else {
			newEnv = append(newEnv, e)
		}
	}
	proc := exec.Command(cmd[0], cmd[1:]...)
	proc.Dir = p.WorkingDirectory
	proc.Env = newEnv
	bufferSize := 1
	stdoutChannel := make(chan string, bufferSize)
	stderrChannel := make(chan string, bufferSize)

	result := []string{}

	stdout, err := proc.StdoutPipe()
	if err != nil {
		return "", RecordError(cmd, err)
	}

	stderr, err := proc.StderrPipe()
	if err != nil {
		return "", RecordError(cmd, err)
	}

	go pipeReader(stdout, stdoutChannel)
	go pipeReader(stderr, stderrChannel)

	logLine := func(line string) {
		if line == "" {
			return
		}
		result = append(result, line)
		log.Log("build.script_output", map[string]string{
			"cmd":  filepath.Base(cmd[0]),
			"line": line,
		})
	}

	done := make(chan int, 1)
	go func() {
		ch1 := stdoutChannel
		ch2 := stderrChannel
		for {
			select {
			case line, ok := <-ch1:
				if !ok {
					ch1 = nil
				} else {
					logLine(line)
				}
			case line, ok := <-ch2:
				if !ok {
					ch2 = nil
				} else {
					logLine(line)
				}
			}
			if ch1 == nil && ch2 == nil {
				done <- 1
				return
			}
		}
	}()

	var returnErr error

	if err := proc.Start(); err != nil {
		returnErr = err
	}
	if err := proc.Wait(); err != nil {
		returnErr = err
	}
	<-done
	lines := strings.Join(result, "\n")
	if returnErr != nil {
		return lines, RecordError(cmd, returnErr)
	}
	return lines, nil
}

func (p *processRecorder) Run(cmd []string, env []string, log api.Logger) error {
	_, err := p.Record(cmd, env, log)
	return err
}

func MakeExecutable(path string) error {
	st, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return os.Chmod(path, st.Mode()|0555)
}
