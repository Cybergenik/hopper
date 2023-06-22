#!/bin/bash

export HOPPER_OUT=$(pwd)"/hopper_out"

go build ./cmd/hopper-node;
./compile examples/parse/getdomain.c

for ((i=$1;i<=$2;i++))
do
    nohup ./hopper-node -I $i -T ./target --raw --args '@@' &> /dev/null &
done
