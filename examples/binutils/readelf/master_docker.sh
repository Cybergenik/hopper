#!/bin/bash

# List: screen -list
# Detach: Ctrl-a Ctrl-d
# Attach: screen -r <pid>.master

## Create Hopper subnet
docker network create hopper-readelf-subnet &> /dev/null

# Master config
HOPPER_OUT="/hopper_out"
CORPUS_PATH="/readelf/corpus"
HAVOC=20

docker run -it --rm \
    --name hopper-master-readelf \
    --env TERM \
    --env HOPPER_OUT=$HOPPER_OUT \
    --env HOPPER_LOG=1 \
    --env HOPPER_LOG_INTERVAL=10 \
    --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
    --volume $(pwd)$CORPUS_PATH:$CORPUS_PATH \
    --network hopper-readelf-subnet \
    --publish 6969:6969 \
    hopper-readelf:latest \
    bash -c "hopper-master -I ${CORPUS_PATH} -H ${HAVOC}"

## Clean up subnet
docker network rm -f hopper-readelf-subnet
