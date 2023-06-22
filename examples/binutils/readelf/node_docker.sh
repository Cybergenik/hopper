#!/bin/bash

## Spawn Nodes
export HOPPER_OUT="/hopper_out"
MASTER_IP="1.1.1.1"
PORT="6969"

for ((i=$1;i<=$2;i++))
do
    nohup docker run --rm \
        --name hopper-node$i \
        --env HOPPER_OUT \
        --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
        hopper-readelf:latest \
        bash -c "
            cd /hopper;
            ./hopper-node -I $i -T ./readelf_target -M $MASTER_IP -P $PORT --args '-a @@'" &> /dev/null &
    sleep .2
done
