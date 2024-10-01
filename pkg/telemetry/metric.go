package telemetry

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewMeter creates a new metric.Meter that can create any metric reporter
// you might want to use in your application.
func NewMeter(ctx context.Context) (metric.Meter, error) {
	provider, err := newMeterProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create meter provider: %w", err)
	}

	return provider.Meter("rest-skeleton"), nil
}

// newMeterProvcider initialize the application resource, connects to the
// OpenTelemetry Collector and configures the metric poller that will be used
// to collect the metrics and send them to the OpenTelemetry Collector.
func newMeterProvider(ctx context.Context) (metric.MeterProvider, error) {
	// Interval which the metrics will be reported to the collector
	interval := 10 * time.Second

	resource, err := getResource()
	if err != nil {
		return nil, fmt.Errorf("could not get resource: %w", err)
	}

	collectorExporter, err := getOtelMetricsCollectorExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get collector exporter: %w", err)
	}

	periodicReader := metricsdk.NewPeriodicReader(collectorExporter,
		metricsdk.WithInterval(interval),
	)

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(resource),
		metricsdk.WithReader(periodicReader),
	)

	return provider, nil
}

// getResource creates the resource that describes our application.
//
// You can add any attributes to your resource and all your metrics
// will contain those attributes automatically.
//
// There are some attributes that are very important to be added to the resource:
// 1. hostname: allows you to identify host-specific problems
// 2. version: allows you to pinpoint problems in specific versions
func getResource() (*resource.Resource, error) {
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(os.Getenv("APP_NAME")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("could not merge resources: %w", err)
	}

	return resource, nil
}

// getOtelMetricsCollectorExporter creates a metric exporter that relies on
// an OpenTelemetry Collector running on "localhost:4317".
func getOtelMetricsCollectorExporter(ctx context.Context) (metricsdk.Exporter, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(os.Getenv("OTEL_COLLECTOR_ENDPOINT")),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create metric exporter: %w", err)
	}

	return exporter, nil
}

func SetMetric(meter metric.Meter) (metric.Int64Histogram, metric.Int64Counter, error) {
	latencyMetric, err := meter.Int64Histogram("http.server.latency")
	if err != nil {
		return nil, nil, fmt.Errorf("could not create metric: %w", err)
	}

	errorCounter, err := meter.Int64Counter("app.error.count")
	if err != nil {
		return nil, nil, fmt.Errorf("could not create metric: %w", err)
	}

	return latencyMetric, errorCounter, nil
}

func CollectMachineResourceMetrics(meter metric.Meter) {
	period := 5 * time.Second
	ticker := time.NewTicker(period)

	var Mb uint64 = 1_048_576 // number of bytes in a MB

	for range ticker.C {
		// Ini akan dieksekusi setiap "period" waktu
		meter.Float64ObservableGauge(
			"process.allocated_memory",
			metric.WithFloat64Callback(
				func(ctx context.Context, fo metric.Float64Observer) error {
					var memStats runtime.MemStats
					runtime.ReadMemStats(&memStats)

					allocatedMemoryInMB := float64(memStats.Alloc) / float64(Mb)
					fo.Observe(allocatedMemoryInMB)

					return nil
				},
			),
		)
	}
}
