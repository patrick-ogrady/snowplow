package utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
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

// Download retrieves a file from a bucket with a
// given name.
func Download(bucket string, name string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("%w: could not create new storage client", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(name).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", name, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	if err := ioutil.WriteFile(name, data, 0644); err != nil {
		return fmt.Errorf("%w: unable to write %s to disk", err, name)
	}

	return nil
}
