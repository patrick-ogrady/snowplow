.PHONY: build install add-license check-license compile docker-build docker-build-local run-mainnet mocks lint test

ADDLICENSE_CMD=go run github.com/google/addlicense
ADDLICENCE_SCRIPT=${ADDLICENSE_CMD} -c "patrick-ogrady" -l "mit" -v
LINT_SETTINGS=golint,misspell,gocyclo,gocritic,whitespace,goconst,gocognit,bodyclose,unconvert,lll,unparam,gomnd
GOLINES_CMD=go run github.com/segmentio/golines

WORKDIR             ?= $(shell pwd)
GIT_COMMIT          ?= $(shell git rev-parse HEAD)
DOCKER_ORG          ?= patrick-ogrady
PROJECT             ?= snowplow
DOCKER_IMAGE        ?= ${DOCKER_ORG}/${PROJECT}
DOCKER_LABEL        ?= latest
DOCKER_TAG          ?= ${DOCKER_IMAGE}:${DOCKER_LABEL}
AVALANCHE_VERSION   ?= v1.1.2
SNOWPLOW_VERSION 		?= v0.0.2

build:
	go build

install:
	go install ./...

add-license:
	${ADDLICENCE_SCRIPT} .;
	go mod tidy;

check-license:
	${ADDLICENCE_SCRIPT} -check .;
	go mod tidy;

docker-build:
	docker build \
		--no-cache \
		--build-arg AVALANCHE_VERSION=${AVALANCHE_VERSION} \
		--build-arg SNOWPLOW_VERSION=${SNOWPLOW_VERSION} \
		-t ${DOCKER_TAG} \
		-f Dockerfile \
		.

docker-build-local:
	docker build \
		--no-cache \
		--build-arg AVALANCHE_VERSION=${AVALANCHE_VERSION} \
		--build-arg SNOWPLOW_VERSION=${GIT_COMMIT} \
		-t ${DOCKER_TAG} \
		-f Dockerfile \
		.

run-mainnet:
	docker run \
		-d \
		-v ${WORKDIR}/.avalanchego:/root/.avalanchego \
		-p 9650:9650 \
		-p 9651:9651 \
		${DOCKER_TAG}

mocks:
	rm -rf mocks;
	mockery --disable-version-string --dir pkg/health --all --case underscore --outpkg health --output mocks/pkg/health;

lint:
	golangci-lint run --timeout 2m0s -v -E ${LINT_SETTINGS}

test:
	go test -v ./pkg/...

shorten-lines:
	${GOLINES_CMD} -w --shorten-comments .;
	go mod tidy;
