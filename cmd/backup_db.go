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
	"github.com/patrick-ogrady/snowplow/pkg/storage"
)

// backupDbCmd represents the backup db command
var backupDbCmd = &cobra.Command{
	Use:   "backup [bucket] [name]",
	Short: "backup db to google cloud storage",
	Args:  cobra.ExactArgs(2), // nolint:gomnd
	RunE:  backupDbFunc,
}

func init() {
	dbCmd.AddCommand(backupDbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func backupDbFunc(cmd *cobra.Command, args []string) error {
	// Check if dbDirectory is empty
	if _, err := os.Stat(dbDirectory); os.IsNotExist(err) {
		return fmt.Errorf("%s is an empty directory", dbDirectory)
	}

	// Tar db
	name := args[1]
	tarFile := fmt.Sprintf("%s.tar.gz", name)
	if err := compression.Compress(dbDirectory, tarFile); err != nil {
		return fmt.Errorf("%w: could not compress db", err)
	}

	// Backup db
	bucket := args[0]
	if err := storage.Upload(
		Context,
		bucket,
		tarFile,
	); err != nil {
		return fmt.Errorf("%w: unable to upload %s", err, tarFile)
	}

	// Cleanup
	if err := os.Remove(tarFile); err != nil {
		return fmt.Errorf("%w: unable to delete %s", err, tarFile)
	}

	fmt.Printf("successfully backed up %s to %s\n", name, bucket)
	return nil
}
