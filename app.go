package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	otlp := newGRPCExporter(context.Background(), "localhost:4317")

	tp, err := tracerProvider(otlp)

	if err != nil {
		panic(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
}

func main() {
	tp := otel.GetTracerProvider()

	mux := http.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {

		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.(*sdktrace.TracerProvider).Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	mux.HandleFunc("/foo", handleManual())

	helloHandler := handleOtelHTTP()
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")

	mux.Handle("/hello", otelHandler)

	log.Fatal(http.ListenAndServe("0.0.0.0:9090", mux))
}

// handleManual instruments HTTP incoming calls manually
func handleManual() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tp := otel.GetTracerProvider()
		t := tp.Tracer("component-http-root-path")
		_, sp := t.Start(r.Context(), "foo")
		sp.SetName("foo")

		sp.AddEvent("Handled with manual instrumentation")

		var attrs []attribute.KeyValue = []attribute.KeyValue{
			attribute.Key("Path").String(r.URL.Path),
		}

		for k, v := range r.Header {
			attrs = append(attrs, attribute.Key(k).String(v[0]))
		}

		sp.SetAttributes(attrs...)
		defer sp.End()

		fmt.Fprint(w, "I am foo.\n")
	}
}

// handleOtelHTTP instruments HTTP incoming calls using go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
func handleOtelHTTP() http.HandlerFunc {
	uk := attribute.Key("username")

	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		bag := baggage.FromContext(ctx)
		span.AddEvent("Handled with otelhttp", trace.WithAttributes(uk.String(bag.Member("username").Value())))

		fmt.Fprint(w, "Hello, world!\n")
	}
}
