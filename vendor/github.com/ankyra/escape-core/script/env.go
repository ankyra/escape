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
	globals[func_builtinSplit] = builtinSplit
	globals[func_builtinJoin] = builtinJoin
	globals[func_builtinBase64Encode] = builtinBase64Encode
	globals[func_builtinBase64Decode] = builtinBase64Decode
	globals[func_builtinReplace] = builtinReplace
	globals[func_builtinTrim] = builtinTrim
	globals[func_builtinAdd] = ShouldLift(builtinAdd)
	globals[func_builtinListIndex] = LiftFunction(builtinListIndex)
	globals[func_builtinListSlice] = LiftFunction(builtinListSlice)
	globals[func_builtinTimestamp] = ShouldLift(builtinTimestamp)
	globals[func_builtinTrackMajorVersion] = ShouldParse(`$func(v) { $v.split(".")[:1].join(".").concat(".@") }`)
	globals[func_builtinTrackMinorVersion] = ShouldParse(`$func(v) { $v.split(".")[:2].join(".").concat(".@") }`)
	globals[func_builtinTrackPatchVersion] = ShouldParse(`$func(v) { $v.split(".")[:3].join(".").concat(".@") }`)
	globals[func_builtinTrackVersion] = ShouldParse(`$func(v) { $v.concat(".@") }`)
	globalsDict := LiftDict(globals)
	result["$"] = globalsDict
	return &result
}
