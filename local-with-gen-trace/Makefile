GOBIN ?= $(shell go env GOPATH)/bin

install-ocb:
	curl --proto '=https' --tlsv1.2 -fL -o ocb https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.113.0/ocb_0.113.0_linux_amd64 && chmod +x ocb

build:
	./ocb --config builder-config.yaml

run:
	./otelcol-dev/otelcol-dev --config ./config.yaml


gen-traces:
	go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@latest\
	&& $(GOBIN)/telemetrygen traces --otlp-insecure --traces 3
