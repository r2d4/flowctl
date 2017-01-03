// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"github.com/golang/glog"

	"github.com/r2d4/flowctl/pkg/flowctl"
	"github.com/spf13/cobra"
)

const defaultNumWorkers = 2
const defaultNumPS = 1

var numWorkers int
var numPS int

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Brings up TensorFlow GRPC servers",
	Long:  `Brings up TensorFlow GRPC servers`,
	Run: func(cmd *cobra.Command, args []string) {
		jobs := map[string]int{
			"worker": 2,
			"ps":     1,
		}
		err := flowctl.CreateServers(jobs)
		if err != nil {
			glog.Exitln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(upCmd)
}
