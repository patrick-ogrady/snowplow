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
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/hashing"
)

// LoadNodeID returns the ID associated with a staker cert.
// Source:
// https://github.com/ava-labs/avalanchego/blob/e2944176f9e87562140ecd979cafebb4707578c4/node/node.go#L407-L430
func LoadNodeID(stakingCertPath string) (ids.ShortID, error) {
	stakeCert, err := ioutil.ReadFile(stakingCertPath)
	if err != nil {
		return ids.ShortID{}, fmt.Errorf("%w: problem reading staking certificate", err)
	}

	block, _ := pem.Decode(stakeCert)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return ids.ShortID{}, fmt.Errorf("%w: problem parsing staking certificate", err)
	}

	id, err := ids.ToShortID(hashing.PubkeyBytesToAddress(cert.Raw))
	if err != nil {
		return ids.ShortID{}, fmt.Errorf("%w: problem deriving staker ID from certificate", err)
	}

	return id, nil
}

// PrintableNodeID returns the canonical form of the NodeID
// for delegation.
func PrintableNodeID(nodeID ids.ShortID) string {
	return nodeID.PrefixedString(constants.NodeIDPrefix)
}
