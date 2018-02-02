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

package validate

var InvalidStateProjectNames = []string{
	"",
	"/",
	"$",
	"a",
	"a/",
	"a$",
	"a^",
	"a*",
	" ",
	"   asddas",
	"@",
	"aa>",
	"<script",
}

var ValidStateProjectNames = []string{
	"aa",
	"ab",
	"test-etsts",
	"test_test",
	"t_____",
	"t------",
	"TEST",
	"1test",
	"1000",
}

var InvalidEnvironmentNames = []string{
	"",
	".../../",
	"$",
	"@",
	":",
	"A",
	"B",
	"1prod",
	"-prod",
	"_prod",
	"PROD",
}
var ValidEnvironmentNames = []string{
	"ci",
	"dev",
	"prod",
	"a",
	"a1",
	"a-1",
	"a-_2",
	"a________3",
}

var InvalidDeploymentNames = []string{
	"",
	"/",
	".",
	",",
	"$",
	"^",
	"-",
	"1test",
	"-test",
	"/test",
}

var ValidDeploymentNames = []string{
	"_",
	"_/test",
	"_/test-test",
	"test",
	"t",
	"TEST",
}
