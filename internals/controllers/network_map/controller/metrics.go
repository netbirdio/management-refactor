package controller

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric"
)

type metrics struct {
	dbAccessDuration metric.Int64Histogram
}

func newMetrics(meter metric.Meter) (*metrics, error) {
	dbAccessDuration, err := meter.Int64Histogram(
		"sync_request_duration_seconds",
		metric.WithDescription("Duration of sync requests in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return &metrics{
		dbAccessDuration: dbAccessDuration,
	}, nil
}

func (m *metrics) RecordDBAccessDuration(duration time.Duration) {
	m.dbAccessDuration.Record(context.Background(), duration.Milliseconds(), metric.WithAttributes())
}
