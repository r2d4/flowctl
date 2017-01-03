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

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Tears down Tensorflow GRPC Servers",
	Long:  `Tears down Tensorflow GRPC Servers`,
	Run: func(cmd *cobra.Command, args []string) {
		err := flowctl.DeleteServers()
		if err != nil {
			glog.Exitln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(downCmd)
}
