#!/bin/bash

for ((i=$1;i<=$2;i++))
do
    nohup docker run --rm \
        -e "TERM=xterm-256color" \
        --network hopper-subnet \
        --name hopper-node$i \
        hopper-node:latest \
        bash -c "
            cd hopper/node;
            go build .;
            cd ..;
            ./compile examples/parse/getdomain.c
            ./node/node -I $i -T ./target -M hopper-master --args '@@'" &> /dev/null &
done
