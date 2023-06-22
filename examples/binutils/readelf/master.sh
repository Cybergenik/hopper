#!/bin/bash

# Logging:
export HOPPER_LOG=1
export HOPPER_LOG_INTERVAL=1

mkdir hopper_out
export HOPPER_OUT=$(pwd)"/hopper_out"
## Spawn Master
go build ./cmd/hopper-master;
./hopper-master -I ./examples/binutils/readelf/in -H=20 #--no-tui

