#!/bin/bash

export HOPPER_OUT=$(pwd)"/hopper_out"

go build .;

./hopper -I ./examples/parse/in -H=5

