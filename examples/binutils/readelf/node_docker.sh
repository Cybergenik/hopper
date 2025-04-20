#!/bin/bash

HOPPER_OUT="/hopper_out"

START_ID=$1
NUM_NODES=$2

## Spawn Nodes
for ((i=START_ID; i < START_ID + NUM_NODES; i++));
do
    nohup docker run --rm \
        --name hopper-node-readelf-$i \
        --env TERM \
        --env HOPPER_OUT=$HOPPER_OUT \
        --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
        --network hopper-readelf-subnet \
        hopper-readelf:latest \
        bash -c "hopper-node -I $i -T /readelf_target -M hopper-master-readelf --args '-a @@'" &> /dev/null &
    echo "Started hopper-node-readelf-${i}"
done
