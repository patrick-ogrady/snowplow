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

// Monitor tracks the health
// of an avalanche validator.
type Monitor struct {
	notifier Notifier
	client   Client

	healthInterval time.Duration
	statusInterval time.Duration

	unhealthyThreshold time.Duration
	minPeers           uint64

	isBootstrappedMutex sync.Mutex
	isBootstrapped      map[string]time.Time

	isHealthy time.Time

	peers    time.Time
	numPeers uint64

	completeHealth            bool
	completeHealthStatusSince time.Time
}

// NewMonitor returns a new *Monitor.
func NewMonitor(
	notifier Notifier,
	client Client,
	healthInterval time.Duration,
	statusInterval time.Duration,
	unhealthyThreshold time.Duration,
	minPeers uint64,
) *Monitor {
	return &Monitor{
		notifier:           notifier,
		client:             client,
		healthInterval:     healthInterval,
		unhealthyThreshold: unhealthyThreshold,
		statusInterval:     statusInterval,
		minPeers:           minPeers,

		isBootstrapped: make(map[string]time.Time),
	}
}

// checkBootstrapped loops on the IsBootstrapped
// check for a particular chain.
func (m *Monitor) checkIsBootstrapped(
	ctx context.Context,
	chain string,
) {
	start := time.Now()
	for utils.ContextSleep(ctx, m.healthInterval) == nil {
		bootstrapped, err := m.client.IsBootstrapped(chain)
		if err != nil {
			m.notifier.Alert(fmt.Sprintf("%s-Chain IsBootstrapped failed: %s", chain, err.Error()))
			continue
		}

		if !bootstrapped {
			continue
		}

		m.notifier.Info(fmt.Sprintf("%s-Chain bootstrapped after %s", chain, time.Since(start)))
		m.isBootstrappedMutex.Lock()
		m.isBootstrapped[chain] = time.Now()
		m.isBootstrappedMutex.Unlock()
		return
	}
}

// checkIsHealthy loops on the IsHealthy
// check.
func (m *Monitor) checkIsHealthy(
	ctx context.Context,
) {
	for utils.ContextSleep(ctx, m.healthInterval) == nil {
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

func (m *Monitor) checkPeers(
	ctx context.Context,
) {
	for utils.ContextSleep(ctx, m.healthInterval) == nil {
		peers, err := m.client.Peers()
		if err != nil {
			m.notifier.Alert(fmt.Sprintf("Peers failed: %s", err.Error()))
			continue
		}

		if peers < m.minPeers {
			continue
		}

		m.peers = time.Now()
	}
}

func (m *Monitor) computeHealth() string {
	m.isBootstrappedMutex.Lock()
	defer m.isBootstrappedMutex.Unlock()
	for _, chain := range chains {
		if _, ok := m.isBootstrapped[chain]; !ok {
			return fmt.Sprintf("%s-Chain isBootstrapped=false", chain)
		}
	}

	if time.Since(m.isHealthy) > m.unhealthyThreshold {
		return fmt.Sprintf("isHealthy=false for %s", time.Since(m.isHealthy))
	}

	if time.Since(m.peers) > m.unhealthyThreshold {
		return fmt.Sprintf("peers < %d for %s", m.minPeers, time.Since(m.peers))
	}

	return ""
}

func (m *Monitor) monitorStatus(ctx context.Context) {
	for utils.ContextSleep(ctx, m.statusInterval) == nil {
		m.notifier.Status(fmt.Sprintf(
			"healthy(%s): %t peers: %d",
			time.Since(m.completeHealthStatusSince),
			m.completeHealth,
			m.numPeers,
		))
	}
}

// MonitorHealth checks a validator's health
// each interval.
func (m *Monitor) MonitorHealth(
	ctx context.Context,
) {
	go m.monitorStatus(ctx)

	for _, chain := range chains {
		go m.checkIsBootstrapped(ctx, chain)
	}
	go m.checkIsHealthy(ctx)
	go m.checkPeers(ctx)

	m.completeHealthStatusSince = time.Now()
	for utils.ContextSleep(ctx, m.healthInterval) == nil {
		healthyStatus := m.computeHealth()

		if (m.completeHealth && len(healthyStatus) == 0) || (!m.completeHealth && len(healthyStatus) > 0) {
			continue
		}

		if m.completeHealth && len(healthyStatus) > 0 {
			m.notifier.Alert(fmt.Sprintf("not healthy: %s", healthyStatus))
			m.completeHealth = false
			m.completeHealthStatusSince = time.Now()
			continue
		}

		if !m.completeHealth && len(healthyStatus) == 0 {
			m.notifier.Info(fmt.Sprintf("healthy after %s", time.Since(m.completeHealthStatusSince)))
			m.completeHealth = true
			m.completeHealthStatusSince = time.Now()
		}
	}
}
