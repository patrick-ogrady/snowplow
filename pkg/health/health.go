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
	"sync"
	"time"

	"github.com/patrick-ogrady/snowplow/pkg/utils"
)

var (
	chains = []string{"X", "C", "P"}
)

// Notifier ...
type Notifier interface {
	Alert(message string)
	Info(message string)
	Status(message string)
}

// Client ...
type Client interface {
	IsHealthy() (bool, error)
	IsBootstrapped(chain string) (bool, error)
	Peers() (uint64, error)
}

// checkBootstrapped loops on the IsBootstrapped
// check for a particular chain.
func (m *Monitor) checkBootstrapped(
	ctx context.Context,
	chain string,
) {
	start := time.Now()
	for utils.ContextSleep(ctx, m.interval) == nil {
		bootstrapped, err := m.client.IsBootstrapped(chain)
		if err != nil {
			m.notifier.Alert(fmt.Sprintf("%s-Chain IsBootstrapped failed: %s", chain, err.Error()))
			continue
		}

		if !bootstrapped {
			continue
		}

		m.notifier.Info(fmt.Sprintf("%s-Chain bootstrapped after %s", chain, time.Since(start)))
		m.bootstrappedMutex.Lock()
		m.bootstrapped[chain] = time.Now()
		m.bootstrappedMutex.Unlock()
		return
	}
}

// checkIsHealthy loops on the IsHealthy
// check.
func (m *Monitor) checkIsHealthy(
	ctx context.Context,
) {
	for utils.ContextSleep(ctx, m.interval) == nil {
		isHealthy, err := m.client.IsHealthy()
		if err != nil {
			m.notifier.Alert(fmt.Sprintf("IsHealthy failed: %s", err.Error()))
			continue
		}

		if !isHealthy {
			continue
		}

		m.isHealthy = time.Now()
	}
}

func (m *Monitor) checkMinPeers(
	ctx context.Context,
) {
	for utils.ContextSleep(ctx, m.interval) == nil {
		peers, err := m.client.Peers()
		if err != nil {
			m.notifier.Alert(fmt.Sprintf("Peers failed: %s", err.Error()))
			continue
		}

		if peers < m.minPeers {
			continue
		}

		m.hasPeers = time.Now()
	}
}

// Monitor tracks the health
// of an avalanche validator.
type Monitor struct {
	interval  time.Duration
	threshold time.Duration
	status    time.Duration
	notifier  Notifier
	client    Client
	minPeers  uint64

	bootstrappedMutex sync.Mutex
	bootstrapped      map[string]time.Time
	isHealthy         time.Time

	hasPeers time.Time
	numPeers uint64

	healthy            bool
	healthyStatusSince time.Time
}

// NewMonitor returns a new *Monitor.
func NewMonitor(
	interval time.Duration,
	threshold time.Duration,
	status time.Duration,
	notifier Notifier,
	client Client,
	minPeers uint64,
) *Monitor {
	return &Monitor{
		interval:  interval,
		threshold: threshold,
		status:    status,
		notifier:  notifier,
		client:    client,
		minPeers:  minPeers,

		bootstrapped: make(map[string]time.Time),
	}
}

func (m *Monitor) checkHealth() bool {
	for _, chain := range chains {
		if _, ok := m.bootstrapped[chain]; !ok {
			return false
		}
	}

	if time.Since(m.isHealthy) > m.threshold {
		return false
	}

	if time.Since(m.hasPeers) > m.threshold {
		return false
	}

	return true
}

func (m *Monitor) checkStatus(ctx context.Context) {
	for utils.ContextSleep(ctx, m.status) == nil {
		m.notifier.Status(fmt.Sprintf(
			"healthy (%s): %t peers: %d",
			time.Since(m.healthyStatusSince),
			m.healthy,
			m.numPeers,
		))
	}
}

// MonitorHealth checks a validator's health
// each interval.
func (m *Monitor) MonitorHealth(
	ctx context.Context,
) {
	for _, chain := range chains {
		go m.checkBootstrapped(ctx, chain)
	}

	go m.checkIsHealthy(ctx)
	go m.checkMinPeers(ctx)
	go m.checkStatus(ctx)

	m.healthyStatusSince = time.Now()
	for utils.ContextSleep(ctx, m.interval) == nil {
		thisHealthy := m.checkHealth()

		if (m.healthy && thisHealthy) || (!m.healthy && !thisHealthy) {
			continue
		}

		if m.healthy && !thisHealthy {
			m.notifier.Alert("not healthy")
			m.healthy = false
			m.healthyStatusSince = time.Now()
			continue
		}

		if !m.healthy && thisHealthy {
			m.notifier.Info(fmt.Sprintf("healthy after %s", time.Since(m.healthyStatusSince)))
			m.healthy = true
			m.healthyStatusSince = time.Now()
		}
	}
}
