package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ava-labs/avalanchego/api/info"
)

const (
	nodeEndpoint = "http://localhost:9560"
	timeout      = time.Second * 10
)

// CheckHealth checks a node's health
// each interval.
func CheckHealth(
	ctx context.Context,
	interval time.Duration,
	notifier *Notifier,
) error {
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

			continue
		}

		if healthy && !thisHealthy.Healthy {
			healthy = false
			notifier.Alert("node no longer healthy")
			continue
		} else if !healthy && thisHealthy.Healthy {
			healthy = true
			notifier.Info("node now healthy")
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

	return ctx.Err()
}
