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

package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/cheggaaa/pb/v3"
	"github.com/machinebox/progress"

	"github.com/patrick-ogrady/snowplow/pkg/integrity"
)

const (
	progressSleepTime = 10 * time.Millisecond
)

// Upload puts a specified file in a bucket with
// a given name.
func Upload(ctx context.Context, bucket string, name string) error {
	// Open local file.
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("%w: could not open file %s", err, name)
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("%w: unable to get file stats for %s", err, name)
	}

	checksum, err := integrity.Checksum(f)
	if err != nil {
		return fmt.Errorf("%w: could not get checksum of credentials", err)
	}

	if err := uploadString(
		ctx,
		bucket,
		fmt.Sprintf("%s.checksum", name),
		checksum,
	); err != nil {
		return fmt.Errorf("%w: unable to upload checksum", err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("%w: unable to reset file pointer", err)
	}

	return upload(ctx, bucket, name, fInfo.Size(), f)
}

// uploadString uploads a string to name.
func uploadString(ctx context.Context, bucket string, name string, blob string) error {
	size := int64(len(blob))
	return upload(ctx, bucket, name, size, bytes.NewReader([]byte(blob)))
}

func logUploadProgress(ctx context.Context, name string, size int64, blob *progress.Reader) {
	go func() {
		fmt.Printf("uploading %s...\n", name)

		bar := pb.Start64(size)
		bar.SetRefreshRate(progressSleepTime)
		for ctx.Err() == nil && (size-blob.N()) > 0 {
			bar.SetCurrent(blob.N())
			time.Sleep(progressSleepTime)
		}
		bar.SetCurrent(size)
		bar.Finish()
	}()
}

func upload(ctx context.Context, bucket string, name string, size int64, blob io.Reader) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("%w: could not create new storage client", err)
	}
	defer client.Close()

	blobProgress := progress.NewReader(blob)
	logUploadProgress(ctx, name, size, blobProgress)

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(name).NewWriter(ctx)
	if _, err = io.Copy(wc, blobProgress); err != nil {
		return fmt.Errorf("%w: io.Copy", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("%w: Writer.Close", err)
	}

	return nil
}

// Download retrieves a file from a bucket with a
// given name.
func Download(ctx context.Context, bucket string, name string) error {
	dChecksum, err := downloadString(
		ctx,
		bucket,
		fmt.Sprintf("%s.checksum", name),
	)
	if err != nil {
		return fmt.Errorf("%w: unable to download checksum", err)
	}

	if err := downloadFile(ctx, bucket, name); err != nil {
		return fmt.Errorf("%w: unable to download %s", err, name)
	}

	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("%w: could not open file %s", err, name)
	}
	defer f.Close()

	checksum, err := integrity.Checksum(f)
	if err != nil {
		return fmt.Errorf("%w: could not get checksum of credentials", err)
	}

	if checksum != dChecksum {
		return fmt.Errorf("expected checksum %s but got %s", dChecksum, checksum)
	}

	return nil
}

// downloadString downloads a string from name.
func downloadString(ctx context.Context, bucket string, name string) (string, error) {
	rc, err := download(ctx, bucket, name)
	if err != nil {
		return "", fmt.Errorf("%w: unable to download %s", err, name)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	return string(data), nil
}

func logDownloadProgress(ctx context.Context, name string, rc *storage.Reader) {
	go func() {
		fmt.Printf("downloading %s...\n", name)

		size := rc.Attrs.Size
		bar := pb.Start64(size)
		bar.SetRefreshRate(progressSleepTime)
		for ctx.Err() == nil && rc.Remain() > 0 {
			bar.SetCurrent(size - rc.Remain())
			time.Sleep(progressSleepTime)
		}
		bar.SetCurrent(size)
		bar.Finish()
	}()
}

// downloadFile downloads a file from name without loading
// it all into memory at once.
func downloadFile(ctx context.Context, bucket string, name string) error {
	rc, err := download(ctx, bucket, name)
	if err != nil {
		return fmt.Errorf("%w: unable to download %s", err, name)
	}
	defer rc.Close()

	w, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("%w: unable to create file %s", err, name)
	}
	defer w.Close()

	if _, err := io.Copy(w, rc); err != nil {
		return fmt.Errorf("%w: unable to download to %s", err, name)
	}

	return nil
}

// download is the core of the download process.
func download(ctx context.Context, bucket string, name string) (*storage.Reader, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: could not create new storage client", err)
	}
	defer client.Close()

	rc, err := client.Bucket(bucket).Object(name).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", name, err)
	}

	logDownloadProgress(ctx, name, rc)
	return rc, nil
}
