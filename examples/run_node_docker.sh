#!/bin/bash

for ((i=$1;i<=$2;i++))
do
    nohup sudo docker run --network=host --rm \
        -e "TERM=xterm-256color" \
        --name hopper-node$i \
        hopper-node:latest \
        bash -c "
            cd hopper/node;
            go build .;
            cd ..;
            ./compile examples/parse/getdomain.c
            ./node/node -I $i -T ./target --args '@@'" &> /dev/null &
done
