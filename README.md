# otel-playground
OpenTelemetry Playground


## Running with in-app exporter directly to Jaeger

  1. Start Jaeger with `docker compose up -d`
  1. Run the application with `go run .`
  1. Opean the Jaeger UI at http://localhost:16686/

## Running with collector exporter via otlp

  1. Start Jaeger with `docker compose up jaeger -d`
  1. Start OTel Collector with `docker compose up collector -d`
  1. Run the application with `go run .`
  1. Opean the Jaeger UI at http://localhost:16686/
