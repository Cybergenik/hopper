#!/bin/bash
set -e

if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Usage: $0 <start_node_id> <num_nodes>"
  exit 1
fi

HOPPER_OUT="/hopper_out"

START_ID=$1
NUM_NODES=$2

## Spawn Nodes
for ((i=START_ID; i < START_ID + NUM_NODES; i++));
do
    nohup docker run --rm \
        --name hopper-node-jq-${i} \
        --env TERM \
        --env HOPPER_OUT=$HOPPER_OUT \
        --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
        --network hopper-jq-subnet \
        hopper-jq:latest \
        bash -c "hopper-node -I $i -T /jq/jq -M hopper-master-jq -P 6969 --stdin" &> /dev/null &

    echo "Started hopper-node-jq-${i}"
done
