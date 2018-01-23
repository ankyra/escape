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

const (
	DEBUG   = iota
	INFO    = iota
	SUCCESS = iota
	WARN    = iota
	ERROR   = iota
)

const EscapeLogo = `
  @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
  @                                 @
  @ @@@@@@                .@.       @
  @ @@@@@@                @@@       @
  @ @@@@@@                @@@.      @
  @ @@@@@@    +             @@@@;   @
  @ @@@@@@   +@            +@@@,@   @
  @ @@@@@@  +@@@@@@@@@@@@  +@@@ :   @
  @ @@@@@@ +@@@@@@@@@@@@@.@@@#@@    @
  @ @@@@@@  +@@@@@@@@@@@@    .@@    @
  @ @@@@@@   +@             #@@@    @
  @ @@@@@@    +             @# @    @ 
  @ @@@@@@                  @, @@@. @
  @ @@@@@@                  @       @
  @ @@@@@@                  @       @
  @                                 @
  @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@`

var logMessages = map[string]map[string]string{
	"package.finished": map[string]string{
		"msg":   "Packaged {{ .release }}",
		"level": "success",
	},
	"package.start": map[string]string{
		"msg":   "Started packaging.",
		"level": "info",
	},
	"build.build_dependency": map[string]string{
		"msg":   "Building dependency {{ .dependency }}.",
		"level": "info",
	},
	"build.build_dependency_finished": map[string]string{
		"msg":   "Built dependency {{ .release }}",
		"level": "success",
	},
	"build.build_step": map[string]string{
		"msg":   "Running build step.",
		"level": "info",
	},
	"build.build_step_finished": map[string]string{
		"msg":   "Successfully ran build step.",
		"level": "success",
	},
	"build.docker": map[string]string{
		"msg":   "Building docker image {{ .image }}.",
		"level": "info",
	},
	"build.docker_no_repository": map[string]string{
		"msg":   "No 'docker_repository' input variable defined. Not pushing.",
		"level": "info",
	},
	"build.docker_push": map[string]string{
		"msg":   "Pushing image {{ .image }}.",
		"level": "info",
	},
	"build.environment_variable_value": map[string]string{
		"msg":   "With environment: {{ .variable }} = {{ .value }}",
		"level": "debug",
	},
	"build.errand_variable_value": map[string]string{
		"msg":   "With extra input errand variable: {{ .variable }} = {{ .value }}",
		"level": "debug",
	},
	"build.finished": map[string]string{
		"msg":   "Completed build",
		"level": "success",
	},
	"build.ignore_output_variable": map[string]string{
		"msg":   "The build output '{{ .variable }} = {{ .value }}' was not defined in the Escape plan. Ignoring.",
		"level": "warn",
	},
	"build.input_variable_value": map[string]string{
		"msg":   "With input: {{ .variable }} = {{ .value }}",
		"level": "debug",
	},
	"build.kubespec": map[string]string{
		"msg":   "Running Kubernetes spec.",
		"level": "info",
	},
	"build.kubespec_apply": map[string]string{
		"msg":   "Applying Kubernetes spec {{ .path }}.",
		"level": "info",
	},
	"build.output_location": map[string]string{
		"msg":   "Outputs (including overrides) can be written to {{ .path }}.",
		"level": "info",
	},
	"build.output_override_variable_value": map[string]string{
		"msg":   "With output override: {{ .variable }} = {{ .value }}",
		"level": "debug",
	},
	"build.output_variable_value": map[string]string{
		"msg":   "With output: {{ .variable }} = {{ .value }}",
		"level": "debug",
	},
	"build.packer": map[string]string{
		"msg":   "Building packer image.",
		"level": "info",
	},
	"build.run_script": map[string]string{
		"msg":   "Running script {{ .cmd }}",
		"level": "debug",
	},
	"build.script_output": map[string]string{
		"msg":      "{{ .cmd }}: {{ .line }}",
		"level":    "info",
		"collapse": "false",
	},
	"build.start": map[string]string{
		"msg":   "Starting build.",
		"level": "info",
	},
	"build.step": map[string]string{
		"msg":   "Running {{ .step }} step {{ .script }}.",
		"level": "info",
	},
	"build.terraform": map[string]string{
		"msg":   "Terraforming.",
		"level": "info",
	},
	"build.terraform_no_outputs_warning": map[string]string{
		"msg":   "Could not get Terraform outputs.",
		"level": "warn",
	},
	"build.terraform_tip": map[string]string{
		"msg":   "TIP: You can create a Terraform file {{ .path }} and Escape will try and parse it.",
		"level": "info",
	},
	"client.add_deployment": map[string]string{
		"msg":   "Deploying {{ .release }} into environment {{ .environment }} of project {{ .project }}.",
		"level": "info",
	},
	"client.authentication_expired": map[string]string{
		"msg":   "Authentication expired. Trying login with stored credentials.",
		"level": "info",
	},
	"client.authentication_failed": map[string]string{
		"msg":   "Authentication failed. Couldn't log in to server: {{ .error_message }}",
		"level": "warn",
	},
	"client.download_release": map[string]string{
		"msg":   "Downloading release {{ .release }}.",
		"level": "info",
	},
	"client.get_environment_state": map[string]string{
		"msg":   "Getting project environment config for {{ .environment }} in project {{ .project }}.",
		"level": "info",
	},
	"client.next_version": map[string]string{
		"msg":   "Querying server for next version of {{ .release }}.",
		"level": "info",
	},
	"client.register": map[string]string{
		"msg":   "Registering release with Escape server.",
		"level": "info",
	},
	"client.release_query": map[string]string{
		"msg":   "Querying release {{ .release }}.",
		"level": "info",
	},
	"client.update_inputs_outputs": map[string]string{
		"msg":   "Sending inputs and outputs for {{ .release }} back to the server.",
		"level": "info",
	},
	"client.upload_release": map[string]string{
		"msg":   "Uploading release to the Escape server.",
		"level": "info",
	},
	"deploy.deploy_dependency": map[string]string{
		"msg":   "Deploying dependency {{ .dependency }}.",
		"level": "info",
	},
	"deploy.deploy_dependency_finished": map[string]string{
		"msg":   "Deployed dependency {{ .release }}",
		"level": "success",
	},
	"deploy.deploy_step": map[string]string{
		"msg":   "Running deployment step.",
		"level": "info",
	},
	"deploy.finished": map[string]string{
		"msg":   "Successfully deployed.",
		"level": "success",
	},
	"deploy.start": map[string]string{
		"msg":   "Deploying.",
		"level": "info",
	},
	"deploy.step": map[string]string{
		"msg":   "Running {{ .step }} step {{ .script }}.",
		"level": "info",
	},
	"deploy.step_finished": map[string]string{
		"msg":   "Successfully ran deployment step.",
		"level": "info",
	},
	"destroy.destroy_dependency": map[string]string{
		"msg":   "Destroying dependency {{ .dependency }}.",
		"level": "info",
	},
	"destroy.destroy_dependency_finished": map[string]string{
		"msg":   "Destroying dependency {{ .release }}",
		"level": "success",
	},
	"destroy.destroy_step": map[string]string{
		"msg":   "Destroying",
		"level": "success",
	},
	"destroy.step_finished": map[string]string{
		"msg":   "Destroy step completed successfully",
		"level": "success",
	},
	"destroy.docker_already_removed": map[string]string{
		"msg":   "Docker image {{ .image }} had already been removed.",
		"level": "success",
	},
	"destroy.docker_finished": map[string]string{
		"msg":   "Successfully removed docker image {{ .image }}.",
		"level": "success",
	},
	"destroy.docker_start": map[string]string{
		"msg":   "Removing docker image {{ .image }}.",
		"level": "info",
	},
	"destroy.finished": map[string]string{
		"msg":   "Destruction complete",
		"level": "success",
	},
	"destroy.nothing_to_do": map[string]string{
		"msg":   "Nothing to destroy.",
		"level": "info",
	},
	"destroy.start": map[string]string{
		"msg":   "Destroying deployment.",
		"level": "info",
	},
	"destroy.terraform_finished": map[string]string{
		"msg":   "Successfully destroyed Terraform resources.",
		"level": "success",
	},
	"destroy.terraform_start": map[string]string{
		"msg":   "Destroying Terraform estate.",
		"level": "info",
	},
	"errand.start": map[string]string{
		"msg":   "Running {{ .errand }}.",
		"level": "info",
	},
	"error": map[string]string{
		"msg":      "Error: {{ .error }}",
		"level":    "error",
		"collapse": "false",
	},
	"fetch.download_from_gcs": map[string]string{
		"msg":   "Downloading {{ .release }} from {{ .gcs_path }} into {{ .target_dir }}.",
		"level": "info",
	},
	"fetch.download_from_gcs_complete": map[string]string{
		"msg":   "Downloaded {{ .release }} from {{ .gcs_path }} into {{ .target_dir }}.",
		"level": "success",
	},
	"fetch.download_from_gcs_failed": map[string]string{
		"msg":      "Failed to download {{ .release }} from {{ .gcs_path }}.",
		"level":    "error",
		"collapse": "false",
	},
	"fetch.finished": map[string]string{
		"msg":   "Fetched dependency {{ .dependency }}",
		"level": "success",
	},
	"fetch.start": map[string]string{
		"msg":   "Fetching dependency {{ .dependency }}.",
		"level": "info",
	},
	"install.finished": map[string]string{
		"msg":   "Finished installing all dependencies.",
		"level": "success",
	},
	"install.start": map[string]string{
		"msg":   "Installing dependencies.",
		"level": "info",
	},
	"login.finished": map[string]string{
		"msg":   "Logged in.",
		"level": "success",
	},
	"plan.written": map[string]string{
		"msg":   "Written {{ .path }}.",
		"level": "success",
	},
	"push.finished": map[string]string{
		"msg":   "Push successful.",
		"level": "success",
	},
	"register.finished": map[string]string{
		"msg":   "Release {{.release}} was successfully registered.",
		"level": "success",
	},
	"register.start": map[string]string{
		"msg":   "Registering.",
		"level": "info",
	},
	"release.finished": map[string]string{
		"msg":   "Successfully released {{.release}}",
		"level": "success",
	},
	"release.start": map[string]string{
		"msg":   "Releasing {{.release}}",
		"level": "info",
	},
	"run.finished": map[string]string{
		"msg":   "Run succeeded.",
		"level": "success",
	},
	"run.start": map[string]string{
		"msg":   "Running.",
		"level": "info",
	},
	"smoke.finished": map[string]string{
		"msg":   "Smoke tests passed.",
		"level": "success",
	},
	"smoke.start": map[string]string{
		"msg":   "Running smoke tests.",
		"level": "info",
	},
	"test.finished": map[string]string{
		"msg":   "Tests passed.",
		"level": "success",
	},
	"test.start": map[string]string{
		"msg":   "Running tests.",
		"level": "info",
	},
	"upload.finished": map[string]string{
		"msg":   "Upload finished.",
		"level": "success",
	},
	"upload.start": map[string]string{
		"msg":   "Uploading release.",
		"level": "info",
	},
	"upload.to_gcs": map[string]string{
		"msg":   "Uploading {{ .path }} to {{ .bucket_url }}",
		"level": "info",
	},
}
