//go:build e2e
// +build e2e

/*
Copyright 2022 The Dapr Authors
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

package standalone_test

import (
	"github.com/dapr/cli/tests/e2e/common"
	"github.com/dapr/cli/tests/e2e/spawn"
)

// cmdDashboard runs the Dapr dashboard and returns the command output and error.
func cmdDashboard(port string) (string, error) {
	return spawn.Command(common.GetDaprPath(), "dashboard", "--port", port)
}

// cmdInit installs Dapr with the init command and returns the command output and error.
//
// When DAPR_E2E_INIT_SLIM is true, it will install Dapr without Docker containers.
// This is useful for scenarios where Docker containers are not available, e.g.,
// in GitHub actions Windows runner.
//
// Arguments to the init command can be passed via args.
func cmdInit(runtimeVersion string, args ...string) (string, error) {
	initArgs := []string{"init", "--log-as-json", "--runtime-version", runtimeVersion}

	if isSlimMode() {
		initArgs = append(initArgs, "--slim")
	}

	initArgs = append(initArgs, args...)

	return spawn.Command(common.GetDaprPath(), initArgs...)
}

// cmdInvoke invokes a method on the specified app and returns the command output and error.
func cmdInvoke(appId, method, unixDomainSocket string, args ...string) (string, error) {
	invokeArgs := []string{"invoke", "--log-as-json", "--app-id", appId, "--method", method}

	if unixDomainSocket != "" {
		invokeArgs = append(invokeArgs, "--unix-domain-socket", unixDomainSocket)
	}

	invokeArgs = append(invokeArgs, args...)

	return spawn.Command(common.GetDaprPath(), invokeArgs...)
}

// cmdList lists the running dapr instances and returns the command output and error.
// format can be empty, "table", "json", or "yaml"
func cmdList(output string) (string, error) {
	args := []string{"list"}

	if output != "" {
		args = append(args, "-o", output)
	}

	return spawn.Command(common.GetDaprPath(), args...)
}

// cmdPublish publishes a message to the specified pubsub and topic, and returns the command output and error.
func cmdPublish(appId, pubsub, topic, unixDomainSocket string, args ...string) (string, error) {
	publishArgs := []string{"publish", "--log-as-json", "--publish-app-id", appId, "--pubsub", pubsub, "--topic", topic}

	if unixDomainSocket != "" {
		publishArgs = append(publishArgs, "--unix-domain-socket", unixDomainSocket)
	}

	publishArgs = append(publishArgs, args...)

	return spawn.Command(common.GetDaprPath(), publishArgs...)
}

// cmdRun runs a Dapr instance and returns the command output and error.
func cmdRun(unixDomainSocket string, args ...string) (string, error) {
	runArgs := append([]string{"run"})

	if unixDomainSocket != "" {
		runArgs = append(runArgs, "--unix-domain-socket", unixDomainSocket)
	}

	runArgs = append(runArgs, args...)

	return spawn.Command(common.GetDaprPath(), runArgs...)
}

// cmdStop stops the specified app and returns the command output and error.
func cmdStop(appId string, args ...string) (string, error) {
	stopArgs := append([]string{"stop", "--log-as-json", "--app-id", appId}, args...)
	return spawn.Command(common.GetDaprPath(), stopArgs...)
}

// cmdUninstall uninstalls Dapr with --all flag and returns the command output and error.
func cmdUninstall() (string, error) {
	return spawn.Command(common.GetDaprPath(), "uninstall", "--all")
}

// cmdVersion checks the version of Dapr and returns the command output and error.
// output can be empty or "json"
func cmdVersion(output string) (string, error) {
	args := []string{"version"}

	if output != "" {
		args = append(args, "-o", output)
	}

	return spawn.Command(common.GetDaprPath(), args...)
}
