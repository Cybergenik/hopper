#!/bin/bash

export HOPPER_OUT="/hopper_out"

for ((i=$1;i<=$2;i++))
do
    nohup docker run --rm \
        --name hopper-node$i \
        --env TERM \
        --env HOPPER_OUT \
        --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
        --network hopper-subnet \
        hopper-node:latest \
        bash -c "
            cd hopper/node;
            go build .;
            cd ..;
            ./compile examples/parse/getdomain.c
            ./node/node -I $i -T ./target -M hopper-master --raw --args '@@'" &> /dev/null &
done
