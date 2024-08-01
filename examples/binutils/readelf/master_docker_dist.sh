#!/bin/bash

# Logging:
export HOPPER_LOG=1
export HOPPER_LOG_INTERVAL=10

# Spawn Master
export HOPPER_OUT="/hopper_out"
LOCAL_OUT="/proj/hopper-tests-PG0/defcon"

docker run -it --rm \
    --name hopper-master \
    --env TERM \
    --env HOPPER_OUT \
    --env HOPPER_LOG \
    --env HOPPER_LOG_INTERVAL \
    --volume $LOCAL_OUT:$HOPPER_OUT \
    --publish 6969:6969 \
    hopper-readelf:latest \
    bash -c "cd hopper && ./hopper-master -I ./examples/binutils/readelf/in -H=20"

