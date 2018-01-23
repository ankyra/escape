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
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"os"
	"reflect"
)

type LogConsumer interface {
	Consume(*LogEntry) error
}

type fancyTerminalOutput struct {
	PreviousSectionStack  []string
	PreviousMessageLength int
	CollapseSections      bool
	LastLineWasCollapsed  bool
}

func NewFancyTerminalOutputLogConsumer() LogConsumer {
	return &fancyTerminalOutput{
		PreviousSectionStack: []string{},
		CollapseSections:     true,
	}
}

func (t *fancyTerminalOutput) Consume(entry *LogEntry) error {

	if entry.Message == "" {
		return nil
	}

	if !t.CollapseSections {
		return t.plainOutput(entry)
	}
	if !entry.Collapse {
		t.PreviousMessageLength = 0
		if t.LastLineWasCollapsed {
			fmt.Println("")
		}
		t.LastLineWasCollapsed = false
		return t.plainOutput(entry)
	}
	switchedCollapsed := !t.LastLineWasCollapsed

	t.LastLineWasCollapsed = true
	sectionChanged := !reflect.DeepEqual(t.PreviousSectionStack, entry.SectionStack)
	if sectionChanged || switchedCollapsed {
		whiteSpace := t.makeWhiteSpace((len(entry.SectionStack) - 1) * 2)
		if !switchedCollapsed {
			fmt.Fprint(os.Stderr, "\n")
		}
		fmt.Fprint(os.Stderr, whiteSpace)
		if len(entry.SectionStack) > 1 {
			fmt.Fprint(os.Stderr, "")
		}
		if len(entry.SectionStack) != 0 {
			fmt.Fprint(os.Stderr, entry.SectionStack[len(entry.SectionStack)-1]+": ")
		}
	} else {
		result := make([]byte, t.PreviousMessageLength)
		space := make([]byte, t.PreviousMessageLength)
		i := 0
		for i < t.PreviousMessageLength {
			result[i] = '\b'
			space[i] = ' '
			i++
		}
		fmt.Fprint(os.Stderr, string(result))
		fmt.Fprint(os.Stderr, string(space))
		fmt.Fprint(os.Stderr, string(result))
	}
	//msg := stripCtlAndExtFromUnicode(entry.Message)
	msg := entry.Message
	t.PreviousMessageLength = len(msg)
	if entry.LogLevel == SUCCESS {
		fmt.Fprint(os.Stderr, "\x1b[32m")
		fmt.Fprint(os.Stderr, "\u2714\ufe0f ")
		t.PreviousMessageLength += 2
	} else if entry.LogLevel == INFO {
	} else if entry.LogLevel == ERROR {
		fmt.Fprint(os.Stderr, "\x1b[31m")
	}
	fmt.Fprint(os.Stderr, msg)
	fmt.Fprint(os.Stderr, "\x1b[0m")
	t.PreviousSectionStack = make([]string, len(entry.SectionStack))
	for i, section := range entry.SectionStack {
		t.PreviousSectionStack[i] = section
	}
	return nil
}

func (t *fancyTerminalOutput) plainOutput(entry *LogEntry) error {
	indent := len(entry.SectionStack) - 1
	if indent < 0 {
		indent = 0
	}
	whiteSpace := t.makeWhiteSpace(indent * 2)
	fmt.Fprint(os.Stderr, whiteSpace)
	if len(entry.SectionStack) > 0 {
		fmt.Fprint(os.Stderr, entry.SectionStack[len(entry.SectionStack)-1])
		fmt.Fprint(os.Stderr, ": ")
	}
	if entry.LogLevel == SUCCESS {
		fmt.Fprint(os.Stderr, "\x1b[32m")
		fmt.Fprint(os.Stderr, "\u2714\ufe0f ")
	} else if entry.LogLevel == ERROR {
		fmt.Fprint(os.Stderr, "\x1b[31m")
	}
	//msg := stripCtlAndExtFromUnicode(entry.Message)
	msg := entry.Message
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprint(os.Stderr, "\x1b[0m")
	return nil
}

func (t *fancyTerminalOutput) makeWhiteSpace(n int) string {
	if n <= 0 {
		return ""
	}
	result := make([]byte, n)
	i := 0
	for i < n {
		result[i] = ' '
		i++
	}
	return string(result)
}

// Advanced Unicode normalization and filtering,
// see http://blog.golang.org/normalization and
// http://godoc.org/golang.org/x/text/unicode/norm for more
// details.
func stripCtlAndExtFromUnicode(str string) string {
	isOk := func(r rune) bool {
		return r < 32 || r >= 127
	}
	// The isOk filter is such that there is no need to chain to norm.NFC
	t := transform.Chain(norm.NFKD, transform.RemoveFunc(isOk))
	// This Transformer could also trivially be applied as an io.Reader
	// or io.Writer filter to automatically do such filtering when reading
	// or writing data anywhere.
	str, _, _ = transform.String(t, str)
	return str
}
