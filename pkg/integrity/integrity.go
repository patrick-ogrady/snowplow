package integrity

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

// Checksum computes a hex-encoded SHA256 checksum of the
// file at filepath.
func Checksum(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("%w: could not open file", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
