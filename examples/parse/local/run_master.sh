#!/bin/bash

# Logging:
#export HOPPER_LOG=1
#export HOPPER_LOG_INTERVAL=1

export HOPPER_OUT=$(pwd)"/hopper_out"

mkdir hopper_out

go build .;

./hopper -I ./examples/parse/in -H=5

