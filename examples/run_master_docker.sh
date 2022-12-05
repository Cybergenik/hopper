#!/bin/bash

# sudo docker build -t hopper-nodes .
## Create Hopper subnet
docker network create hopper-subnet &> /dev/null
## Spawn Master
docker run -it --rm \
    -e "TERM=xterm-256color" \
    --name hopper-master \
    --network hopper-subnet \
    --publish 6969:6969 \
    hopper-node:latest \
    bash -c "cd hopper; go build .; ./hopper -I ./examples/parse/in -H=5 && cat hopper.report"

