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

package util

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ProcessRecorder interface {
	SetWorkingDirectory(string)
	Record(cmd []string, env []string, log Logger) (string, error)
	Run(cmd []string, env []string, log Logger) error
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

func (p *processRecorder) Record(cmd []string, env []string, log Logger) (string, error) {
	proc := exec.Command(cmd[0], cmd[1:]...)
	proc.Dir = p.WorkingDirectory
	proc.Env = env
	bufferSize := 2
	quitChannel := make(chan int)
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
	scanner := bufio.NewScanner(stdout)
	go func(resultChannel chan string) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			resultChannel <- line
		}
		close(resultChannel)
	}(stdoutChannel)
	errScanner := bufio.NewScanner(stderr)
	go func(resultChannel chan string) {
		for errScanner.Scan() {
			line := strings.TrimSpace(errScanner.Text())
			resultChannel <- line
		}
		close(resultChannel)
	}(stderrChannel)

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

	go func(stdoutChannel, stderrChannel chan string, quit chan int) {
		for {
			select {
			case line := <-stdoutChannel:
				logLine(line)
			case line := <-stderrChannel:
				logLine(line)
			case <-quit:
				return
			}
		}
	}(stdoutChannel, stderrChannel, quitChannel)

	if err := proc.Start(); err != nil {
		quitChannel <- 0
		drainChannels(stdoutChannel, stderrChannel, logLine)
		return strings.Join(result, "\n"), RecordError(cmd, err)
	}
	if err := proc.Wait(); err != nil {
		quitChannel <- 0
		drainChannels(stdoutChannel, stderrChannel, logLine)
		return strings.Join(result, "\n"), RecordError(cmd, err)
	}
	quitChannel <- 0
	drainChannels(stdoutChannel, stderrChannel, logLine)
	return strings.Join(result, "\n"), nil
}

func drainChannels(stdoutChannel, stderrChannel chan string, logLine func(string)) {
	for line := range stdoutChannel {
		logLine(line)
	}
	for line := range stderrChannel {
		logLine(line)
	}
}

func (p *processRecorder) Run(cmd []string, env []string, log Logger) error {
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
