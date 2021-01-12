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

	"github.com/spf13/cobra"

	"github.com/patrick-ogrady/snowplow/pkg/compression"
	"github.com/patrick-ogrady/snowplow/pkg/encryption"
	"github.com/patrick-ogrady/snowplow/pkg/storage"
	"github.com/patrick-ogrady/snowplow/pkg/utils"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore [bucket] [node ID]",
	Short: "restore staking credentials from google cloud storage",
	RunE:  restoreFunc,
	Args:  cobra.ExactArgs(2), // nolint:gomnd
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func restoreFunc(cmd *cobra.Command, args []string) error {
	// Check if stakingDirectory is empty
	if _, err := os.Stat(stakingDirectory); !os.IsNotExist(err) {
		return fmt.Errorf("%s is not empty directory", stakingDirectory)
	}

	// Download Credentials
	bucket := args[0]
	printableNodeID := args[1]
	encryptedFilePath := fmt.Sprintf("%s.zip.gpg", printableNodeID)
	if err := storage.Download(
		Context,
		bucket,
		encryptedFilePath,
	); err != nil {
		return fmt.Errorf("%w: unable to download %s", err, encryptedFilePath)
	}

	// Decrypt
	zipFile := fmt.Sprintf("%s.zip", printableNodeID)
	if err := encryption.Decrypt(encryptedFilePath, zipFile); err != nil {
		return fmt.Errorf("%w: could not decrypt credentials", err)
	}

	// Unzip
	if err := compression.Decompress(zipFile, "."); err != nil {
		return fmt.Errorf("%w: could not unzip %s", err, zipFile)
	}

	// Verify Credential Matches
	nodeID, err := utils.LoadNodeID(stakingCertPath)
	if err != nil {
		return fmt.Errorf("%w: could not calculate recovered NodeID", err)
	}
	recoveredNodeID := utils.PrintableNodeID(nodeID)
	if printableNodeID != recoveredNodeID {
		return fmt.Errorf(
			"recovered NodeID %s does not match requested NodeID %s",
			recoveredNodeID,
			printableNodeID,
		)
	}

	// Cleanup
	if err := os.Remove(zipFile); err != nil {
		return fmt.Errorf("%w: unable to delete %s", err, zipFile)
	}
	if err := os.Remove(encryptedFilePath); err != nil {
		return fmt.Errorf("%w: unable to delete %s", err, encryptedFilePath)
	}

	fmt.Printf("successfully restored %s to %s\n", printableNodeID, stakingDirectory)
	return nil
}
