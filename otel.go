package main

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func tracerProvider(exp sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	// use the Jaeger exporter if one is not provided
	var err error

	// Jeager URL: "http://localhost:14268/api/traces"

	if exp == nil {
		// jaeger.WithEndpoint("")
		exp, err = jaeger.New(jaeger.WithCollectorEndpoint())

		if err != nil {
			return nil, err
		}
	}

	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("otel playground"),
			attribute.String("environment", "development"),
			attribute.Int64("ID", int64(os.Getpid())),
		)),
	)

	return tp, nil
}

func newGRPCExporter(ctx context.Context, endpoint string, additionalOpts ...otlptracegrpc.Option) *otlptrace.Exporter {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithReconnectionPeriod(50 * time.Millisecond),
	}

	opts = append(opts, additionalOpts...)
	client := otlptracegrpc.NewClient(opts...)
	exp, _ := otlptrace.New(ctx, client)

	return exp
}
