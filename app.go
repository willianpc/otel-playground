package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	ctx := context.Background()

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp, err := tracerProvider("http://localhost:14268/api/traces")

	if err != nil {
		panic(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		provider := otel.GetTracerProvider()

		t := provider.Tracer("component-http-root-path")

		_, sp := t.Start(r.Context(), "http")

		var attrs []attribute.KeyValue = []attribute.KeyValue{
			attribute.Key("Path").String(r.URL.Path),
		}

		for k, v := range r.Header {
			attrs = append(attrs, attribute.Key(k).String(v[0]))
		}

		sp.SetAttributes(attrs...)
		defer sp.End()

		fmt.Fprint(w, "hello!")
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:9090", nil))
}
