package health

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "github.com/patrick-ogrady/avalanche-runner/mocks/pkg/health"
)

func handleBootstrappedChecks(t *testing.T, n *mocks.Notifier, c *mocks.Client, chain string) {
	c.On("IsBootstrapped", chain).Return(false, nil).Once()
	c.On("IsBootstrapped", chain).Return(false, errors.New("bad")).Once()
	n.On("Alert", fmt.Sprintf("%s-Chain IsBootstrapped check failed: bad", chain)).Once()
	c.On("IsBootstrapped", chain).Return(true, nil).Once()
	n.On("Info", mock.Anything).Run(
		func(args mock.Arguments) {
			// We cannot check explicit chains here because we use
			// mock.Anything as the argument.
			assert.Contains(t, args[0], "-Chain bootstrapped after")
		},
	).Once()
}

func handleHealthChecks(t *testing.T, cancel context.CancelFunc, n *mocks.Notifier, c *mocks.Client) {
	c.On("IsHealthy").Return(false, nil).Once()
	c.On("IsHealthy").Return(false, nil).Once()
	c.On("IsHealthy").Return(false, nil).Once()
	c.On("IsHealthy").Return(true, nil).Once()
	n.On("Info", mock.Anything).Run(
		func(args mock.Arguments) {
			assert.Contains(t, args[0], "healthy after")
		},
	).Once()
	c.On("IsHealthy").Return(false, errors.New("unable to complete health check")).Once()
	n.On("Alert", mock.Anything).Run(
		func(args mock.Arguments) {
			assert.Contains(t, args[0], "IsHealthy check failed: unable to complete health check")
		},
	).Once()
	c.On("IsHealthy").Return(true, nil).Once()
	n.On("Info", mock.Anything).Run(
		func(args mock.Arguments) {
			assert.Contains(t, args[0], "healthy after")
		},
	).Once()
	c.On("IsHealthy").Return(false, nil).Once()
	n.On("Alert", mock.Anything).Run(
		func(args mock.Arguments) {
			assert.Contains(t, args[0], "not healthy")
		},
	).Once()
	c.On("IsHealthy").Return(true, nil).Once()
	n.On("Info", mock.Anything).Run(
		func(args mock.Arguments) {
			assert.Contains(t, args[0], "healthy after")
		},
	).Once()
	c.On("IsHealthy").Return(true, nil).Run(
		func(args mock.Arguments) {
			cancel()
		},
	).Once()
}

func TestMonitorHealth(t *testing.T) {
	notifier := &mocks.Notifier{}
	client := &mocks.Client{}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	handleBootstrappedChecks(t, notifier, client, "X")
	handleBootstrappedChecks(t, notifier, client, "C")
	handleBootstrappedChecks(t, notifier, client, "P")

	handleHealthChecks(t, cancel, notifier, client)

	MonitorHealth(ctx, 10*time.Millisecond, notifier, client)

	time.Sleep(500 * time.Millisecond)
	cancel()

	client.AssertExpectations(t)
	notifier.AssertExpectations(t)
}
