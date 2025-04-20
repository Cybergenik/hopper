#!/bin/bash

# Master config
HOPPER_OUT="/hopper_out"
CORPUS_PATH="/corpus"
HAVOC=20

docker run -it --rm \
    --name hopper-master-readelf \
    --env TERM \
    --env HOPPER_OUT=$HOPPER_OUT \
    --env HOPPER_LOG=1 \
    --env HOPPER_LOG_INTERVAL=10 \
    --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
    --volume $(pwd)$CORPUS_PATH:$CORPUS_PATH \
    --publish 6969:6969 \
    hopper-readelf:latest \
    bash -c "hopper-master -I ${CORPUS_PATH} -H ${HAVOC}"
