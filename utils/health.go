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
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ava-labs/avalanchego/api/info"
)

const (
	nodeEndpoint = "http://localhost:9650"
	timeout      = time.Second * 10
)

// CheckHealth checks a node's health
// each interval.
func CheckHealth(
	ctx context.Context,
	nodeID string,
	interval time.Duration,
) {
	notifier, err := NewNotifier(nodeID)
	if err != nil {
		fmt.Printf("not initializing notifier: %s\n", err.Error())
		return
	}

	healthClient := health.NewClient(nodeEndpoint, timeout)
	infoClient := info.NewClient(nodeEndpoint, timeout)

	var healthy, bootstrapped bool
	for ctx.Err() != nil {
		time.Sleep(interval)

		thisHealthy, err := healthClient.GetLiveness()
		if err != nil {
			if healthy {
				notifier.Alert(fmt.Sprintf("health check failed: %s", err.Error()))
			}

			fmt.Printf("received error while checking liveness: %s\n", err.Error())
			continue
		}

		if healthy && !thisHealthy.Healthy {
			healthy = false
			notifier.Alert("node no longer healthy")
			continue
		} else if !healthy && thisHealthy.Healthy {
			healthy = true
			notifier.Info("node now healthy")
		} else if !healthy && !thisHealthy.Healthy {
			fmt.Println("node not yet healthy")
			continue
		}

		if bootstrapped {
			continue
		}

		xBootstrapped, err := infoClient.IsBootstrapped("X")
		if err != nil {
			notifier.Alert(fmt.Sprintf("X-Chain bootstap check failed: %s", err.Error()))
			continue
		}

		pBootstrapped, err := infoClient.IsBootstrapped("P")
		if err != nil {
			notifier.Alert(fmt.Sprintf("P-Chain bootstap check failed: %s", err.Error()))
			continue
		}

		cBootstrapped, err := infoClient.IsBootstrapped("C")
		if err != nil {
			notifier.Alert(fmt.Sprintf("C-Chain bootstap check failed: %s", err.Error()))
			continue
		}

		if !bootstrapped && xBootstrapped && pBootstrapped && cBootstrapped {
			bootstrapped = true
			notifier.Info("all chains bootstapped")
		}
	}
}
