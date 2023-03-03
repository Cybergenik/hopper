#!/bin/bash

## Spawn Nodes
export HOPPER_OUT="/proj/hopper-tests-PG0/readelf-dat"
MASTER_IP="10.10.1.1"
PORT="6969"

for ((i=$1;i<=$2;i++))
do
    nohup ./node/node -I $i -T ./readelf_target -M $MASTER_IP -P $PORT --args '-a @@' &> /dev/null &
done
