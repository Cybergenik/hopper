#!/bin/bash

# List: screen -list
# Detach: Ctrl-a Ctrl-d
# Attach: screen -r <pid>.master
## Spawn Master
export HOPPER_OUT="/proj/hopper-tests-PG0/readelf-dat"
screen -S master -d -m docker run -it --rm \
    --name hopper-master \
    --env TERM \
    --env HOPPER_OUT \
    --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
    --publish 6969:6969 \
    hopper-readelf:latest \
    bash -c "cd hopper; go build .; ./hopper -I ./examples/binutils/readelf/in -H=5"
