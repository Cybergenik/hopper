#!/bin/bash
set -e

if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Usage: $0 <start_node_id> <num_nodes>"
  exit 1
fi

HOPPER_OUT="/hopper_out"
TARGET="/getdomain.c"

START_ID=$1
NUM_NODES=$2

## Spawn Nodes
for ((i=START_ID; i < START_ID + NUM_NODES; i++));
do
    nohup docker run --rm \
        --name hopper-node-parse-${i} \
        --env TERM \
        --env HOPPER_OUT=$HOPPER_OUT \
        --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
        --volume $(pwd)$TARGET:$TARGET \
        --network hopper-parse-subnet \
        hopper-node:latest \
        bash -c "asan_compile /getdomain.c &&
            hopper-node -I $i -T ./target -M hopper-master-parse --raw --args '@@'" &> /dev/null &

    echo "Started hopper-parse-node-${i}"
done
