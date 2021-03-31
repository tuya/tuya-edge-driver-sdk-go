.PHONY: build test clean docker

GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=device-service-template/cmd/device-service-template
.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
DOCKER_TAG=$(VERSION)-dev

GOFLAGS=-ldflags "-X github.com/tuya/tuya-edge-driver-sdk-go.Version=$(VERSION)"
GOTESTFLAGS?=-race

GIT_SHA=$(shell git rev-parse HEAD)

build: $(MICROSERVICES)
	$(GO) install -tags=safe

device-service-template/cmd:
	$(GO) build $(GOFLAGS) -o $@ ./device-service-template/cmd/device-service-template

docker:
	docker build \
		-f device-service-template/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t tuya/device-service-template:$(GIT_SHA) \
		-t tuya/device-service-template:$(DOCKER_TAG) \
		.

test:
	GO111MODULE=on go test $(GOTESTFLAGS) -coverprofile=coverage.out ./...
	GO111MODULE=on go vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	./bin/test-attribution-txt.sh
	./bin/test-go-mod-tidy.sh

clean:
	rm -f $(MICROSERVICES)
