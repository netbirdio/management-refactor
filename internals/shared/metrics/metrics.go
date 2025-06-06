package metrics

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/exporters/prometheus"
	metric2 "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const defaultEndpoint = "/metrics"

type AppMetrics struct {
	meter    metric2.Meter
	listener net.Listener
	ctx      context.Context
}

func NewAppMetrics() (*AppMetrics, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	pkg := reflect.TypeOf(defaultEndpoint).PkgPath()
	meter := provider.Meter(pkg)

	return &AppMetrics{
		meter: meter,
	}, nil
}

// Expose metrics on a given port and endpoint. If endpoint is empty a defaultEndpoint one will be used.
// Exposes metrics in the Prometheus format https://prometheus.io/
func (appMetrics *AppMetrics) Expose(ctx context.Context, port int, endpoint string) error {
	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	rootRouter := mux.NewRouter()
	rootRouter.Handle(endpoint, promhttp.HandlerFor(
		prometheus2.DefaultGatherer,
		promhttp.HandlerOpts{EnableOpenMetrics: true}))
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	appMetrics.listener = listener
	go func() {
		err := http.Serve(listener, rootRouter)
		if err != nil {
			return
		}
	}()

	log.WithContext(ctx).Infof("enabled application metrics and exposing on http://%s", listener.Addr().String())

	return nil
}

// Close stop application metrics HTTP handler and closes listener.
func (appMetrics *AppMetrics) Close() error {
	if appMetrics.listener == nil {
		return nil
	}
	return appMetrics.listener.Close()
}

// func (appMetrics *AppMetrics) RegisterMetrics(fn func(meter metric2.Meter) error) error {
// 	return fn(appMetrics.meter)
// }

func RegisterMetrics[T any](app *AppMetrics, fn func(metric2.Meter) (T, error)) (T, error) {
	return fn(app.meter)
}
