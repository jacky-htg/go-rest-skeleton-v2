package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0" // Update to the latest version
)

func InitTracing() (func(context.Context) error, error) {
	// Create a gRPC connection
	conn, err := grpc.Dial(os.Getenv("OTEL_COLLECTOR_ENDPOINT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	// Initialize the OTLP trace exporter with gRPC
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithGRPCConn(conn),
			otlptracegrpc.WithCompressor("gzip"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	// Create a new tracer provider with batching
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(os.Getenv("APP_NAME")),
			),
		),
	)

	// Set the global tracer provider
	otel.SetTracerProvider(tracerProvider)

	// Return a function to shutdown the tracer provider
	return tracerProvider.Shutdown, nil
}
