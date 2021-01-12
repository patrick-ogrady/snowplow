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

package health

import (
	"context"
	"fmt"
	"time"
)

// Notifier ...
type Notifier interface {
	Alert(message string)
	Info(message string)
}

// Client ...
type Client interface {
	IsHealthy() (bool, error)
	IsBootstrapped(chain string) (bool, error)
}

// MonitorHealth checks a node's health
// each interval.
func MonitorHealth(
	ctx context.Context,
	interval time.Duration,
	notifier Notifier,
	client Client,
) {
	var healthy, bootstrapped, sentHealth, sentBootstrapped bool
	startHealth := time.Now()

	for ctx.Err() == nil {
		time.Sleep(interval)

		thisHealthy, err := client.IsHealthy()
		if err != nil {
			if healthy {
				notifier.Alert(fmt.Sprintf("health check failed: %s", err.Error()))
			}

			fmt.Printf("received error while checking health: %s\n", err.Error())
			continue
		}

		if healthy && !thisHealthy {
			healthy = false
			notifier.Alert("not healthy")
			continue
		} else if !healthy && thisHealthy {
			healthy = true
			if !sentHealth {
				sentHealth = true
				notifier.Info(fmt.Sprintf("healthy after %s", time.Since(startHealth)))
			} else {
				notifier.Info("healthy")
			}
		} else if !healthy && !thisHealthy {
			continue
		}

		if bootstrapped {
			continue
		}

		xBootstrapped, err := client.IsBootstrapped("X")
		if err != nil {
			notifier.Alert(fmt.Sprintf("X-Chain bootstap check failed: %s", err.Error()))
			continue
		}

		pBootstrapped, err := client.IsBootstrapped("P")
		if err != nil {
			notifier.Alert(fmt.Sprintf("P-Chain bootstap check failed: %s", err.Error()))
			continue
		}

		cBootstrapped, err := client.IsBootstrapped("C")
		if err != nil {
			notifier.Alert(fmt.Sprintf("C-Chain bootstap check failed: %s", err.Error()))
			continue
		}

		if !bootstrapped && xBootstrapped && pBootstrapped && cBootstrapped {
			bootstrapped = true
			if !sentBootstrapped {
				sentBootstrapped = true
				notifier.Info(fmt.Sprintf("chains bootstrapped after %s", time.Since(startHealth)))
			} else {
				notifier.Info("chains bootstapped")
			}
			continue
		}

		fmt.Printf(
			"chains not yet bootstrapped: c-chain=%t x-chain=%t p-chain=%t\n",
			cBootstrapped,
			xBootstrapped,
			pBootstrapped,
		)
	}
}