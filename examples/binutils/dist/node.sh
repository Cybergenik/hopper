#!/bin/bash

## Spawn Nodes
export HOPPER_OUT="/hopper_out"
LOCAL_OUT="/proj/hopper-tests-PG0/readelf-dat"
MASTER_IP="node0"
PORT="6969"

for ((i=$1;i<=$2;i++))
do
    nohup docker run --rm \
        --name hopper-node$i \
        --env HOPPER_OUT \
        --volume $LOCAL_OUT:$HOPPER_OUT \
        hopper-readelf:latest \
        bash -c "
            cd hopper/node;
            go build .;
            cd ..;
            ./node/node -I $i -T ./readelf_target -M $MASTER_IP -P $PORT --args '-a @@'" &> /dev/null &
done
