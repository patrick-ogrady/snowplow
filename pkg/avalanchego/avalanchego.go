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

package avalanchego

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/patrick-ogrady/snowplow/pkg/client"
	"github.com/patrick-ogrady/snowplow/pkg/health"
	"github.com/patrick-ogrady/snowplow/pkg/notifier"
)

const (
	avalanchegoBin  = "/app/avalanchego"
	avalancheConfig = "/app/avalanchego-config.json"

	healthCheckInterval = time.Second * 10
)

// Run starts an avalanchego node.
func Run(ctx context.Context, nodeID string, notifier *notifier.Notifier) error {
	cmd := exec.Command(
		avalanchegoBin,
		"--config-file",
		avalancheConfig,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Send interrupt signal if context is
	// done
	go func() {
		<-ctx.Done()
		if cmd.Process != nil {
			_ = cmd.Process.Signal(os.Interrupt)
		}
	}()

	// Periodically check health and send
	// notifications as needed
	go health.MonitorHealth(
		ctx,
		healthCheckInterval,
		notifier,
		client.NewClient(),
	)

	return cmd.Run()
}
