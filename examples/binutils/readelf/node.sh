#!/bin/bash

set -e

export HOPPER_OUT=$(pwd)"/hopper_out"

# Build readelf target
if [ ! -f ./readelf_target ]
then
    ./examples/binutils/build-binutils.py -t arm-linux-gnueabi
    cp ./examples/binutils/install/bin/arm-linux-gnueabi-readelf ./readelf_target
    rm -rf ./examples/binutils/binutils-* ./examples/binutils/install ./examples/binutils/build
fi


# build Node
go build ./cmd/hopper-node

# Create Nodes
for ((i=$1;i<=$2;i++))
do
    nohup ./hopper-node -I $i -T ./readelf_target -M localhost --args '-a @@' &> /dev/null &
done
