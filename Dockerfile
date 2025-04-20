FROM golang:1.23-alpine AS builder

WORKDIR /hopper
COPY . .

RUN go build -o hopper-master ./cmd/hopper-master \
    && go build -o hopper-node   ./cmd/hopper-node

FROM debian:stable-slim

COPY --from=builder /hopper/hopper-master /usr/local/bin/
COPY --from=builder /hopper/hopper-node   /usr/local/bin/
COPY --from=builder /hopper/asan_compile /usr/local/bin/

## Install clang utils
RUN apt-get update && apt-get install -y --no-install-recommends \
    libclang-dev llvm-dev lld libclang-rt-dev \
    clang libasan8 clang-tools gcc g++ make \
    wget git ca-certificates \
    autoconf automake libtool m4 pkg-config build-essential \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 6969/tcp
