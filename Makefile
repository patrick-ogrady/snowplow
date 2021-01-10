.PHONY: build install add-license check-license compile docker-build docker-build-local run-mainnet

ADDLICENSE_CMD=go run github.com/google/addlicense
ADDLICENCE_SCRIPT=${ADDLICENSE_CMD} -c "patrick-ogrady" -l "mit" -v

GIT_COMMIT          ?= $(shell git rev-parse HEAD)
WORKDIR             ?= $(shell pwd)
DOCKER_ORG          ?= patrick-ogrady
PROJECT             ?= avalanche-runner
DOCKER_IMAGE        ?= ${DOCKER_ORG}/${PROJECT}
DOCKER_LABEL        ?= latest
DOCKER_TAG          ?= ${DOCKER_IMAGE}:${DOCKER_LABEL}
AVALANCHE_VERSION   ?= v1.1.2
RUNNER_VERSION 			?= v0.0.1

build:
	go build ./...

install:
	go install ./...

add-license:
	${ADDLICENCE_SCRIPT} .;
	go mod tidy;

check-license:
	${ADDLICENCE_SCRIPT} -check .;
	go mod tidy;

compile:
	./scripts/compile.sh $(version)

docker-build:
	docker build \
		--no-cache \
		--build-arg AVALANCHE_VERSION=${AVALANCHE_VERSION} \
		--build-arg RUNNER_VERSION=${RUNNER_VERSION} \
		-t ${DOCKER_TAG} \
		-f Dockerfile \
		.

docker-build-local:
	docker build \
		--no-cache \
		--build-arg AVALANCHE_VERSION=${AVALANCHE_VERSION} \
		--build-arg RUNNER_VERSION=${GIT_COMMIT} \
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
