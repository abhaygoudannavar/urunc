// Copyright (c) 2023-2026, Nubificus LTD
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

// containerState represents the platform agnostic pieces relating to a
// running container's status and state, matching the runc implementation.
type containerState struct {
	Version     string            `json:"ociVersion"`
	ID          string            `json:"id"`
	Pid         int               `json:"pid"`
	Status      string            `json:"status"`
	Bundle      string            `json:"bundle"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

var stateCommand = &cli.Command{
	Name:  "state",
	Usage: "output the state of a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container.

EXAMPLE:
For example, if the container id is "ubuntu01" the following will output the
state of "ubuntu01":

	# urunc state ubuntu01`,
	Description: `The state command outputs current state information for the
container instance, conforming to the OCI runtime specification.`,
	Action: func(_ context.Context, cmd *cli.Command) error {
		logrus.WithField("command", "STATE").WithField("args", os.Args).Debug("urunc INVOKED")
		if err := checkArgs(cmd, 1, exactArgs); err != nil {
			return err
		}

		// get Unikontainer data from state.json
		unikontainer, err := getUnikontainer(cmd)
		if err != nil {
			return err
		}

		status := unikontainer.State.Status
		pid := unikontainer.State.Pid

		// Dynamically check if the process is still running.
		// If it has died, OCI spec requires pid to be 0 and status to be stopped.
		if pid > 0 {
			if err := syscall.Kill(pid, 0); err != nil {
				status = "stopped"
				pid = 0
			}
		}

		cs := containerState{
			Version:     unikontainer.State.Version,
			ID:          unikontainer.State.ID,
			Pid:         pid,
			Status:      status,
			Bundle:      unikontainer.State.Bundle,
			Annotations: unikontainer.State.Annotations,
		}

		stateData, err := json.MarshalIndent(cs, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal container state: %w", err)
		}

		// OCI runtime spec requires state to be printed to stdout
		_, err = fmt.Fprintln(os.Stdout, string(stateData))
		return err
	},
}
