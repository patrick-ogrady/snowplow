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

package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// Encrypt encrypts a file using gpg.
func Encrypt(input string, output string) error {
	encryptCmd := exec.Command(
		"gpg",
		"--symmetric",
		"--cipher-algo",
		"aes256",
		"--digest-algo",
		"sha256",
		"--cert-digest-algo",
		"sha256",
		"--compress-algo",
		"none",
		"-z",
		"0",
		"--s2k-mode",
		"3",
		"--s2k-digest-algo",
		"sha512",
		"--s2k-count",
		"65011712",
		"--force-mdc",
		"--no-symkey-cache",
		"-o",
		output,
		"-c",
		input,
	)
	encryptCmd.Stdin = os.Stdin
	encryptCmd.Stdout = os.Stdout
	encryptCmd.Stderr = os.Stderr
	if err := encryptCmd.Run(); err != nil {
		return fmt.Errorf("%w: could not encrypt %s", err, input)
	}

	return nil
}

// Decrypt decrypts a file using gpg.
func Decrypt(input string, output string) error {
	decryptCmd := exec.Command(
		"gpg",
		"--no-symkey-cache",
		"-o",
		output,
		"--decrypt",
		input,
	)
	decryptCmd.Stdin = os.Stdin
	decryptCmd.Stdout = os.Stdout
	decryptCmd.Stderr = os.Stderr
	if err := decryptCmd.Run(); err != nil {
		return fmt.Errorf("%w: could not decrypt %s", err, input)
	}

	return nil
}
