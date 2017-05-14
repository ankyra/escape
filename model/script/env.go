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

package script

import ()

type ScriptEnvironment map[string]Script

func NewScriptEnvironment() *ScriptEnvironment {
	result := ScriptEnvironment{}
	return &result
}
func NewScriptEnvironmentFromMap(m map[string]Script) *ScriptEnvironment {
	result := ScriptEnvironment(m)
	return &result
}

func NewScriptEnvironmentWithGlobals(globals map[string]Script) *ScriptEnvironment {
	result := ScriptEnvironment{}
	if globals == nil {
		globals = map[string]Script{}
	}
	globals[func_builtinId] = LiftFunction(builtinId)
	globals[func_builtinEnvLookup] = LiftFunction(builtinEnvLookup)
	globals[func_builtinConcat] = LiftFunction(builtinConcat)
	globals[func_builtinToLower] = builtinToLower
	globals[func_builtinToUpper] = builtinToUpper
	globals[func_builtinTitle] = builtinTitle
	globalsDict := LiftDict(globals)
	result["$"] = globalsDict
	return &result
}
