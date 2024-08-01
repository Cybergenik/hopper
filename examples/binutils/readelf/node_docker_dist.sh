#!/bin/bash

export HOPPER_OUT="/hopper_out"
LOCAL_OUT="/proj/hopper-tests-PG0/defcon"
MASTER_IP="10.10.1.1"
PORT="6969"

## Spawn Nodes
for ((i=$1;i<=$2;i++))
do
    nohup docker run --rm \
        --name hopper-node$i \
        --env TERM \
        --env HOPPER_OUT \
        --volume $LOCAL_OUT:$HOPPER_OUT \
        hopper-readelf:latest \
        bash -c "
            cd /hopper;
            ./hopper-node -I $i -T ./readelf_target -M $MASTER_IP -P $PORT --args '-a @@'" &> /dev/null &
done
