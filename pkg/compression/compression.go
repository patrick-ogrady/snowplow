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

package compression

import (
	"fmt"
	"os/exec"
)

// Compress compresses a file using tar.
func Compress(input string, output string) error {
	tarCmd := exec.Command(
		"tar",
		"-cvzf",
		output,
		input,
	)
	if err := tarCmd.Run(); err != nil {
		return fmt.Errorf("%w: could not compress %s", err, input)
	}

	return nil
}

// Decompress decompresses a file using tar.
func Decompress(input string, output string) error {
	untarCmd := exec.Command(
		"tar",
		"-xvf",
		input,
		"-C",
		output,
	)
	if err := untarCmd.Run(); err != nil {
		return fmt.Errorf("%w: could not uncompress %s", err, input)
	}

	return nil
}
