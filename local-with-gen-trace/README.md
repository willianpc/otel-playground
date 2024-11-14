### Install ocb

    curl --proto '=https' --tlsv1.2 -fL -o ocb \
    https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.113.0/ocb_0.113.0_linux_amd64

    chmod +x ocb

Or run `make install-ocb`

### Create a builder-config.yaml

To be used to build the collector. All config for the collector should be there, will be used by ocb.

### Build the collector

    ./ocb --config builder-config.yaml

Or run `make build`

### Run the collector

    ./otelcol-dev/otelcol-dev --config ./config.yaml

Or run `make run`

### Simulate traces

`telemetrygen` will generate fake traces to be sent to the collector.
Run this in a separate terminal. Then check stdout (make sure debugexporter is running) or filename.json (if fileexporter is running).


    export GOBIN=${GOBIN:-$(go env GOPATH)/bin}

    go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@latest

    $GOBIN/telemetrygen traces --otlp-insecure --traces 3

Or run `make gen-traces`
