VERSION ?= $(shell cat version.txt | tr -d '\n')

.PHONE: sync-model
sync-model:
	rm -rf model
	curl https://codeload.github.com/open-telemetry/semantic-conventions/tar.gz/main | tar -xz --strip=1 semantic-conventions-main/model

.PHONY: build
build:
	@echo "Building hny-otel-semantic $(VERSION)"
	go build -ldflags="-X 'main.VERSION=$(VERSION)'" -o hny-otel-semantic
