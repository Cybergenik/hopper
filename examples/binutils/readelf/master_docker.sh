#!/bin/bash

# List: screen -list
# Detach: Ctrl-a Ctrl-d
# Attach: screen -r <pid>.master

## Logger:
#export HOPPER_LOG=1
#export HOPPER_LOG_INTERVAL=10

## Create Hopper subnet
docker network create hopper-subnet &> /dev/null

## Spawn Master
export HOPPER_OUT="/hopper_out"

docker run -it --rm \
    --name hopper-master \
    --env TERM \
    --env HOPPER_OUT \
    --env HOPPER_LOG \
    --env HOPPER_LOG_INTERVAL \
    --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
    --network hopper-subnet \
    --publish 6969:6969 \
    hopper-readelf:latest \
    bash -c "cd hopper && ./hopper-master -I ./examples/binutils/readelf/in -H=20"

## Clean up subnet
docker network rm hopper-subnet &> /dev/null
