// Copyright (c) 2021 patrick-ogrady
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/patrick-ogrady/avalanche-runner/pkg/avalanchego"
	"github.com/patrick-ogrady/avalanche-runner/pkg/utils"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run avalanchego",
	RunE:  runFunc,
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runFunc(cmd *cobra.Command, args []string) error {
	// Modify staking paths (avalanchego binary assumes all in homeDir)
	stakingDirectory = filepath.Join(homeDir, stakingDirectory)
	stakingKeyPath = filepath.Join(homeDir, stakingKeyPath)
	stakingCertPath = filepath.Join(homeDir, stakingCertPath)

	// Check if stakingDirectory is empty
	if _, err := os.Stat(stakingDirectory); os.IsNotExist(err) {
		return fmt.Errorf("%s is an empty directory", stakingDirectory)
	}

	// Check if staking key exists
	if _, err := os.Stat(stakingKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("staking key at %s does not exist", stakingKeyPath)
	}

	// Check if staking certificate exists
	if _, err := os.Stat(stakingCertPath); os.IsNotExist(err) {
		return fmt.Errorf("staking certificate at %s does not exist", stakingCertPath)
	}

	nodeID, err := utils.LoadNodeID(stakingCertPath)
	if err != nil {
		return fmt.Errorf("%w: could not calculate NodeID", err)
	}

	// Run avalanchego
	return avalanchego.Run(Context, utils.PrintableNodeID(nodeID))
}
