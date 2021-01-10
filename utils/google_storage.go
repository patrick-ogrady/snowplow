package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// Upload puts a specified file in a bucket with
// a given name.
func Upload(bucket string, name string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("%w: could not create new storage client", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("%w: could not open file", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(name).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("%w: io.Copy", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("%w: Writer.Close", err)
	}

	return nil
}

// func Download() error {
// }
