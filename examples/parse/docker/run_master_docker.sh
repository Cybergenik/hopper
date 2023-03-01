#!/bin/bash

# sudo docker build -t hopper-node .
## Create Hopper subnet
docker network create hopper-subnet &> /dev/null
export HOPPER_OUT="/hopper_out"
## Spawn Master
docker run -it --rm \
    --name hopper-master \
    --env TERM \
    --env HOPPER_OUT \
    --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
    --network hopper-subnet \
    --publish 6969:6969 \
    hopper-node:latest \
    bash -c "cd hopper; go build .; ./hopper -I ./examples/parse/in -H=5"
## Clean up subnet
docker network rm hopper-subnet &> /dev/null

