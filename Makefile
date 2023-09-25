
.PHONE: sync-model
sync-model:
	rm -rf model
	curl https://codeload.github.com/open-telemetry/semantic-conventions/tar.gz/main | tar -xz --strip=1 semantic-conventions-main/model

.PHONY: build
build: sync-model
	go build -o hny-otel-semantic

