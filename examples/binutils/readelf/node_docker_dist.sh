#!/bin/bash

HOPPER_OUT="/hopper_out"
LOCAL_OUT="/proj/hopper-tests-PG0/defcon"
MASTER_IP="10.10.1.1"
PORT="6969"

START_ID=$1
NUM_NODES=$2

## Spawn Nodes
for ((i=START_ID; i < START_ID + NUM_NODES; i++));
do
    nohup docker run --rm \
        --name hopper-node-readelf-$i \
        --env TERM \
        --env HOPPER_OUT=$HOPPER_OUT \
        --volume $LOCAL_OUT:$HOPPER_OUT \
        hopper-readelf:latest \
        bash -c "hopper-node -I $i -T readelf_target -M $MASTER_IP -P $PORT --args '-a @@'" &> /dev/null &
    echo "Started hopper-node-readelf-${i}"
done
