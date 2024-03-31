#!/bin/bash

export HOPPER_OUT="/hopper_out"

## Spawn Nodes
for ((i=$1;i<=$2;i++))
do
    nohup docker run --rm \
        --name hopper-node$i \
        --env TERM \
        --env HOPPER_OUT \
        --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
        --network hopper-subnet \
        hopper-readelf:latest \
        bash -c "
            cd /hopper;
            ./hopper-node -I $i -T ./readelf_target -M hopper-master --args '-a @@'" &> /dev/null &
done
