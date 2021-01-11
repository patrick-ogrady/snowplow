# Copyright (c) 2021 patrick-ogrady
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

# Inspired By: https://github.com/figment-networks/avalanche-rosetta

# ------------------------------------------------------------------------------
# Build avalanche
# ------------------------------------------------------------------------------
FROM golang:1.15 AS avalanche

ARG AVALANCHE_VERSION

RUN git clone https://github.com/ava-labs/avalanchego.git \
  /go/src/github.com/ava-labs/avalanchego

WORKDIR /go/src/github.com/ava-labs/avalanchego

RUN git checkout $AVALANCHE_VERSION && \
    ./scripts/build.sh

# ------------------------------------------------------------------------------
# Build avalanche runner
# ------------------------------------------------------------------------------
FROM golang:1.15 AS runner

ARG RUNNER_VERSION

RUN git clone https://github.com/patrick-ogrady/avalanche-runner.git \
  /go/src/github.com/patrick-ogrady/avalanche-runner

WORKDIR /go/src/github.com/patrick-ogrady/avalanche-runner

ENV CGO_ENABLED=1
ENV GOARCH=amd64
ENV GOOS=linux

RUN git checkout $RUNNER_VERSION && \
    go mod download

RUN \
  GO_VERSION=$(go version | awk {'print $3'}) \
  GIT_COMMIT=$(git rev-parse HEAD) \
  make build

# ------------------------------------------------------------------------------
# Target container for running the node
# ------------------------------------------------------------------------------
FROM ubuntu:18.04

# Install dependencies
RUN apt-get update -y && \
    apt-get install -y wget

WORKDIR /app

# Install avalanche binaries
COPY --from=avalanche \
  /go/src/github.com/ava-labs/avalanchego/build/avalanchego \
  /app/avalanchego

# Install plugins
COPY --from=avalanche \
  /go/src/github.com/ava-labs/avalanchego/build/plugins/* \
  /app/plugins/

# Install avalanche runner
COPY --from=runner \
  /go/src/github.com/patrick-ogrady/avalanche-runner/avalanche-runner \
  /app/avalanche-runner

# Install config
COPY --from=runner \
  /go/src/github.com/patrick-ogrady/avalanche-runner/assets/avalanchego-config.json \
  /app/avalanchego-config.json

EXPOSE 9650
EXPOSE 9651

ENTRYPOINT ["/app/avalanche-runner", "run"]
