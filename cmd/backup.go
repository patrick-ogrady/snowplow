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
	"github.com/patrick-ogrady/snowplow/pkg/integrity"
	"github.com/patrick-ogrady/snowplow/pkg/storage"
	"github.com/patrick-ogrady/snowplow/pkg/utils"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup [bucket]",
	Short: "backup staking credentials to google cloud storage",
	Args:  cobra.ExactArgs(1),
	RunE:  backupFunc,
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func backupFunc(cmd *cobra.Command, args []string) error {
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

	// Load NodeID
	nodeID, err := utils.LoadNodeID(stakingCertPath)
	if err != nil {
		return fmt.Errorf("%w: could not calculate NodeID", err)
	}
	printableNodeID := utils.PrintableNodeID(nodeID)

	// Tar Credentials
	tarFile := fmt.Sprintf("%s.tar.gz", printableNodeID)
	if err := compression.Compress(stakingDirectory, tarFile); err != nil {
		return fmt.Errorf("%w: could not compress credentials", err)
	}

	// Encrypt Credentials
	encryptedFilePath := fmt.Sprintf("%s.gpg", tarFile)
	if err := encryption.Encrypt(tarFile, encryptedFilePath); err != nil {
		return fmt.Errorf("%w: could not encrypt credentials", err)
	}

	// Checksum credentials
	checksum, err := integrity.Checksum(encryptedFilePath)
	if err != nil {
		return fmt.Errorf("%w: could not get checksum of credentials", err)
	}

	// Upload checksum
	bucket := args[0]
	if err := storage.UploadString(
		Context,
		bucket,
		fmt.Sprintf("%s.checksum", printableNodeID),
		checksum,
	); err != nil {
		return fmt.Errorf("%w: unable to upload checksum", err)
	}

	// Backup Credentials
	if err := storage.Upload(
		Context,
		bucket,
		encryptedFilePath,
	); err != nil {
		return fmt.Errorf("%w: unable to upload %s", err, encryptedFilePath)
	}

	// Cleanup
	if err := os.Remove(tarFile); err != nil {
		return fmt.Errorf("%w: unable to delete %s", err, tarFile)
	}
	if err := os.Remove(encryptedFilePath); err != nil {
		return fmt.Errorf("%w: unable to delete %s", err, encryptedFilePath)
	}

	fmt.Printf("successfully backed up %s to %s\n", printableNodeID, bucket)
	return nil
}
