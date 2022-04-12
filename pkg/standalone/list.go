/*
Copyright 2021 The Dapr Authors
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

package standalone

import (
	"strconv"
	"strings"
	"time"

	ps "github.com/mitchellh/go-ps"
	process "github.com/shirou/gopsutil/process"

	"github.com/dapr/cli/pkg/age"
	"github.com/dapr/cli/pkg/metadata"
	"github.com/dapr/cli/utils"
	"github.com/dapr/dapr/pkg/runtime"
)

// ListOutput represents the application ID, application port and creation time.
type ListOutput struct {
	AppID          string `csv:"APP ID"    json:"appId"          yaml:"appId"`
	HTTPPort       int    `csv:"HTTP PORT" json:"httpPort"       yaml:"httpPort"`
	GRPCPort       int    `csv:"GRPC PORT" json:"grpcPort"       yaml:"grpcPort"`
	AppPort        int    `csv:"APP PORT"  json:"appPort"        yaml:"appPort"`
	MetricsEnabled bool   `csv:"-"         json:"metricsEnabled" yaml:"metricsEnabled"` // Not displayed in table, consumed by dashboard.
	Command        string `csv:"COMMAND"   json:"command"        yaml:"command"`
	Age            string `csv:"AGE"       json:"age"            yaml:"age"`
	Created        string `csv:"CREATED"   json:"created"        yaml:"created"`
	CliPID         int    `csv:"CLI PID"   json:"cliPid"         yaml:"cliPid"`
	DaprdPID       int    `csv:"DAPRD PID" json:"daprdPid"       yaml:"daprdPid"`
}

func (d *daprProcess) List() ([]ListOutput, error) {
	return List()
}

// List outputs all the applications.
func List() ([]ListOutput, error) {
	list := []ListOutput{}

	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}

	// Populates the list if all data is available for the sidecar.
	for _, proc := range processes {
		executable := strings.ToLower(proc.Executable())
		if (executable == "daprd") || (executable == "daprd.exe") {
			procDetails, err := process.NewProcess(int32(proc.Pid()))
			if err != nil {
				continue
			}

			cmdLine, err := procDetails.Cmdline()
			if err != nil {
				continue
			}

			cmdLineItems := strings.Fields(cmdLine)
			if len(cmdLineItems) <= 1 {
				continue
			}

			argumentsMap := make(map[string]string)
			for i := 1; i < len(cmdLineItems)-1; i += 2 {
				argumentsMap[cmdLineItems[i]] = cmdLineItems[i+1]
			}

			httpPort := runtime.DefaultDaprHTTPPort
			if httpPortArg, ok := argumentsMap["--dapr-http-port"]; ok {
				if httpPortArgInt, err := strconv.Atoi(httpPortArg); err == nil {
					httpPort = httpPortArgInt
				}
			}

			grpcPort := runtime.DefaultDaprAPIGRPCPort
			if grpcPortArg, ok := argumentsMap["--dapr-grpc-port"]; ok {
				if grpcPortArgInt, err := strconv.Atoi(grpcPortArg); err == nil {
					grpcPort = grpcPortArgInt
				}
			}

			appPort, err := strconv.Atoi(argumentsMap["--app-port"])
			if err != nil {
				appPort = 0
			}

			enableMetrics, err := strconv.ParseBool(argumentsMap["--enable-metrics"])
			if err != nil {
				// Default is true for metrics.
				enableMetrics = true
			}
			appID := argumentsMap["--app-id"]
			appCmd := ""
			cliPIDString := ""
			socket := argumentsMap["--unix-domain-socket"]
			appMetadata, err := metadata.Get(httpPort, appID, socket)
			if err == nil {
				appCmd = appMetadata.Extended["appCommand"]
				cliPIDString = appMetadata.Extended["cliPID"]
			}

			// Parse functions return an error on bad input.
			cliPID, err := strconv.Atoi(cliPIDString)
			if err != nil {
				cliPID = 0
			}

			daprdPid := proc.Pid()

			createUnixTimeMilliseconds, err := procDetails.CreateTime()
			if err != nil {
				continue
			}

			createTime := time.Unix(createUnixTimeMilliseconds/1000, 0)

			listRow := ListOutput{
				Created: createTime.Format("2006-01-02 15:04.05"),
				Age:     age.GetAge(createTime),
				CliPID:  cliPID,
			}

			listRow.AppID = appID
			listRow.HTTPPort = httpPort
			listRow.GRPCPort = grpcPort
			listRow.AppPort = appPort
			listRow.MetricsEnabled = enableMetrics
			listRow.Command = utils.TruncateString(appCmd, 20)
			listRow.DaprdPID = daprdPid

			// filter only dashboard instance.
			if listRow.AppID != "" {
				list = append(list, listRow)
			}
		}
	}

	return list, nil
}
