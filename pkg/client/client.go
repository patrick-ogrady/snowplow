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

package client

import (
	"time"

	"github.com/ava-labs/avalanchego/utils/rpc"
)

const (
	nodeEndpoint = "http://localhost:9650"
	timeout      = time.Second * 10
)

// Client for Avalanche API Endpoints
// TODO: migrate to official client once
// `json: cannot unmarshal object into Go struct field Result.checks.error of
// type error` is fixed.
type Client struct {
	healthRequester rpc.EndpointRequester
	infoRequester   rpc.EndpointRequester
}

// NewClient ...
// Inspired by: https://github.com/ava-labs/avalanchego/blob/8be88a342fced5522cd503b72f49aae450eea863/api/health/client.go
func NewClient() *Client {
	return &Client{
		healthRequester: rpc.NewEndpointRequester(nodeEndpoint, "/ext/health", "health", timeout),
		infoRequester:   rpc.NewEndpointRequester(nodeEndpoint, "/ext/info", "info", timeout),
	}
}

// GetLivenessReply is the response for GetLiveness
type GetLivenessReply struct {
	Checks  map[string]interface{} `json:"checks"`
	Healthy bool                   `json:"healthy"`
}

// GetLiveness returns a health check on the Avalanche node
func (c *Client) GetLiveness() (*GetLivenessReply, error) {
	res := &GetLivenessReply{}
	err := c.healthRequester.SendRequest("getLiveness", struct{}{}, res)
	return res, err
}

// IsHealthy ...
func (c *Client) IsHealthy() (bool, error) {
	liveness, err := c.GetLiveness()
	if err != nil {
		return false, err
	}

	return liveness.Healthy, nil
}

// IsBootstrappedArgs are the arguments for calling IsBootstrapped
type IsBootstrappedArgs struct {
	// Alias of the chain
	// Can also be the string representation of the chain's ID
	Chain string `json:"chain"`
}

// IsBootstrappedResponse are the results from calling IsBootstrapped
type IsBootstrappedResponse struct {
	// True iff the chain exists and is done bootstrapping
	IsBootstrapped bool `json:"isBootstrapped"`
}

// IsBootstrapped ...
func (c *Client) IsBootstrapped(chain string) (bool, error) {
	res := &IsBootstrappedResponse{}
	err := c.infoRequester.SendRequest("isBootstrapped", &IsBootstrappedArgs{
		Chain: chain,
	}, res)
	return res.IsBootstrapped, err
}
