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

package metrics

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredres "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

const (
	instanceURL = "http://metadata.google.internal/computeMetadata/v1/instance/id"
	zoneURL     = "http://metadata.google.internal/computeMetadata/v1/instance/zone"
	projectURL  = "http://metadata.google.internal/computeMetadata/v1/project/project-id"

	peersMetric  = "custom.googleapis.com/peers"
	minMetricGap = 10 * time.Second
)

// MetricWriter writes metrics to Google Cloud Monitoring.
type MetricWriter struct {
	client *monitoring.MetricClient

	projectID  string
	instanceID string
	zone       string

	nodeID string

	lastWrite      map[string]time.Time
	lastWriteMutex sync.Mutex
}

func loadStringAttribute(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

// NewMetricWriter creates a new *MetricWriter.
func NewMetricWriter(ctx context.Context, nodeID string) (*MetricWriter, error) {
	project, err := loadStringAttribute(projectURL)
	if err != nil {
		return nil, fmt.Errorf("%w: could not load project id", err)
	}

	extendedZone, err := loadStringAttribute(zoneURL)
	if err != nil {
		return nil, fmt.Errorf("%w: could not load zone", err)
	}
	zoneComponents := strings.Split(extendedZone, "/")
	zone := zoneComponents[len(zoneComponents)-1]

	instance, err := loadStringAttribute(instanceURL)
	if err != nil {
		return nil, fmt.Errorf("%w: could not load instance id", err)
	}

	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: unable to create metric client", err)
	}

	return &MetricWriter{
		client:     client,
		projectID:  project,
		zone:       zone,
		instanceID: instance,
		nodeID:     nodeID,
		lastWrite:  make(map[string]time.Time),
	}, nil
}

// Close closes all connections held by the *MetricWriter.
func (w *MetricWriter) Close() error {
	return w.client.Close()
}

func (w *MetricWriter) writeInt64(ctx context.Context, metric string, num int64) error {
	w.lastWriteMutex.Lock()
	defer w.lastWriteMutex.Unlock()

	v, ok := w.lastWrite[metric]
	if ok && time.Since(v) < minMetricGap {
		return nil
	}

	now := &timestamp.Timestamp{
		Seconds: time.Now().Unix(),
	}
	req := &monitoringpb.CreateTimeSeriesRequest{
		Name: "projects/" + w.projectID,
		TimeSeries: []*monitoringpb.TimeSeries{{
			Metric: &metricpb.Metric{
				Type: metric,
				Labels: map[string]string{
					"nodeID": w.nodeID,
				},
			},
			Resource: &monitoredres.MonitoredResource{
				Type: "gce_instance",
				Labels: map[string]string{
					"instance_id": w.instanceID,
					"zone":        w.zone,
				},
			},
			Points: []*monitoringpb.Point{{
				Interval: &monitoringpb.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_Int64Value{
						Int64Value: num,
					},
				},
			}},
		}},
	}

	if err := w.client.CreateTimeSeries(ctx, req); err != nil {
		return fmt.Errorf("%w: could not write time series value", err)
	}

	w.lastWrite[metric] = time.Now()
	return nil
}

// Peers writes the peerCount to metrics.
func (w *MetricWriter) Peers(ctx context.Context, peerCount uint64) error {
	if w == nil {
		return nil
	}

	return w.writeInt64(ctx, peersMetric, int64(peerCount))
}
