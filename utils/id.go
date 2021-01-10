package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/hashing"
)

// NodeID returns the ID associated with a staker cert.
// Source: https://github.com/ava-labs/avalanchego/blob/e2944176f9e87562140ecd979cafebb4707578c4/node/node.go#L407-L430
func NodeID(stakingCertPath string) (ids.ShortID, error) {
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
