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
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "github.com/patrick-ogrady/snowplow/mocks/pkg/health"
)

func handleIsBootstrappedChecks(t *testing.T, n *mocks.Notifier, c *mocks.Client, chain string) {
	c.On("IsBootstrapped", chain).Return(false, nil).Once()
	c.On("IsBootstrapped", chain).Return(false, errors.New("bad")).Once()
	n.On("Alert", fmt.Sprintf("%s-Chain IsBootstrapped failed: bad", chain)).Once()
	c.On("IsBootstrapped", chain).Return(true, nil).Once()
	n.On("Info", mock.Anything).Run(
		func(args mock.Arguments) {
			// We cannot check explicit chains here because we use
			// mock.Anything as the argument.
			assert.Contains(t, args[0], "-Chain bootstrapped after")
		},
	).Once()
}

func handleIsHealthyChecks(n *mocks.Notifier, c *mocks.Client) {
	c.On("IsHealthy").Return(false, nil).Once()
	c.On("IsHealthy").Return(true, nil).Once()
	c.On("IsHealthy").Return(false, errors.New("unable to complete health check")).Once()
	n.On("Alert", "IsHealthy failed: unable to complete health check").Once()
	c.On("IsHealthy").Return(true, nil).Once()
	// should not send a healthy recovery because of threshold
	c.On("IsHealthy").Return(true, nil).Once()
	c.On("IsHealthy").Return(false, nil).Once()
	c.On("IsHealthy").Return(false, nil).Once()
	c.On("IsHealthy").Return(false, nil).Once()
}

func handlePeers(n *mocks.Notifier, c *mocks.Client) {
	c.On("Peers").Return(uint64(0), nil).Once()
	c.On("Peers").Return(uint64(2), nil).Once()
	c.On("Peers").Return(uint64(3), nil).Once()
	c.On("Peers").Return(uint64(4), nil).Once()
	c.On("Peers").Return(uint64(5), nil).Once()
	n.On("Info", "connected peers (5) >= 5").Once()
	c.On("Peers").Return(uint64(5), nil).Once()
	c.On("Peers").Return(uint64(5), nil).Once()
	c.On("Peers").Return(uint64(5), nil).Once()
}

func handleStatus(ctx context.Context, t *testing.T, n *mocks.Notifier) {
	var seenTrue bool
	n.On("Status", mock.Anything).Run(
		func(args mock.Arguments) {
			if strings.Contains(args[0].(string), "true") {
				seenTrue = true
			}
		},
	)

	go func() {
		<-ctx.Done()
		assert.True(t, seenTrue)
	}()
}

func TestMonitorHealth(t *testing.T) {
	notifier := &mocks.Notifier{}
	client := &mocks.Client{}
	metricWriter := &mocks.MetricWriter{}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	for _, chain := range chains {
		handleIsBootstrappedChecks(t, notifier, client, chain)
	}
	handlePeers(notifier, client)
	handleIsHealthyChecks(notifier, client)
	handleStatus(ctx, t, notifier)

	notifier.On("Info", mock.Anything).Run(
		func(args mock.Arguments) {
			assert.Contains(t, args[0], "healthy after")
		},
	).Once()

	notifier.On("Alert", mock.Anything).Run(
		func(args mock.Arguments) {
			assert.Contains(t, args[0], "not healthy: isHealthy=false for")
			cancel()
		},
	).Once()

	m := NewMonitor(notifier, client, metricWriter, 10*time.Millisecond, 15*time.Millisecond, 30*time.Millisecond, 5)
	m.MonitorHealth(ctx)

	time.Sleep(500 * time.Millisecond)
	client.AssertExpectations(t)
	notifier.AssertExpectations(t)
	metricWriter.AssertExpectations(t)
}
