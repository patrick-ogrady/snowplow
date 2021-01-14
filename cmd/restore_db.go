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

// restoreDbCmd represents the restore db command
var restoreDbCmd = &cobra.Command{
	Use:   "db [bucket] [name]",
	Short: "restore db from google cloud storage",
	RunE:  restoreDbFunc,
	Args:  cobra.ExactArgs(2), // nolint:gomnd
}

func init() {
	restoreCmd.AddCommand(restoreDbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func restoreDbFunc(cmd *cobra.Command, args []string) error {
	// Check if dbDirectory is empty
	if _, err := os.Stat(dbDirectory); !os.IsNotExist(err) {
		return fmt.Errorf("%s is not empty directory", dbDirectory)
	}

	// Download backup
	bucket := args[0]
	name := args[1]
	tarFilePath := fmt.Sprintf("%s.tar.gz", name)
	if err := storage.Download(
		Context,
		bucket,
		tarFilePath,
	); err != nil {
		return fmt.Errorf("%w: unable to download %s", err, tarFilePath)
	}

	// Untar credentials
	if err := compression.Decompress(tarFilePath, "."); err != nil {
		return fmt.Errorf("%w: could not decompress %s", err, tarFilePath)
	}

	// Cleanup
	if err := os.Remove(tarFilePath); err != nil {
		return fmt.Errorf("%w: unable to delete %s", err, tarFilePath)
	}

	fmt.Printf("successfully restored %s to %s\n", name, dbDirectory)
	return nil
}
