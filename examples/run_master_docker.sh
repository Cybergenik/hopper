#!/bin/bash

# sudo docker build -t hopper-nodes .
sudo docker run -it --rm \
    -e "TERM=xterm-256color" \
    --name hopper-master \
    -p 6969:6969/tcp \
    hopper-node:latest \
    bash -c "cd hopper; go build .; ./hopper -I ./examples/parse/in -H=5"

