package metrics

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
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
func NewMetricWriter(nodeID string) (*MetricWriter, error) {
	project, err := loadStringAttribute(projectURL)
	if err != nil {
		return nil, fmt.Errorf("%w: could not load project id", err)
	}

	zone, err := loadStringAttribute(zoneURL)
	if err != nil {
		return nil, fmt.Errorf("%w: could not load zone", err)
	}

	instance, err := loadStringAttribute(instanceURL)
	if err != nil {
		return nil, fmt.Errorf("%w: could not load instance id", err)
	}

	return &MetricWriter{
		projectID:  project,
		zone:       zone,
		instanceID: instance,
	}, nil
}

func (w *MetricWriter) writeInt64(ctx context.Context, metric string, num int64) error {
	w.lastWriteMutex.Lock()
	defer w.lastWriteMutex.Unlock()

	v, ok := w.lastWrite[metric]
	if ok && time.Since(v) > minMetricGap {
		return nil
	}

	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
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
					"nodeId": w.nodeID,
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

	if err := c.CreateTimeSeries(ctx, req); err != nil {
		return fmt.Errorf("%w: could not write time series value", err)
	}

	w.lastWrite[metric] = time.Now()
	return nil
}

// Peers writes the peerCount to metrics.
func (w *MetricWriter) Peers(ctx context.Context, peerCount uint64) error {
	return w.writeInt64(ctx, peersMetric, int64(peerCount))
}
