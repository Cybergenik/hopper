#!/bin/bash

set -xe
# Logging:
#export HOPPER_LOG=1
#export HOPPER_LOG_INTERVAL=1

export HOPPER_OUT=$(pwd)"/hopper_out"

mkdir hopper_out

go build ./cmd/hopper-master;

./hopper-master -I ./examples/parse/in -H=5

