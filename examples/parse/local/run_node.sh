#!/bin/bash

export HOPPER_OUT=$(pwd)"/hopper_out"

cd node;
go build .;
cd ..;
./compile examples/parse/getdomain.c

for ((i=$1;i<=$2;i++))
do
    nohup ./node/node -I $i -T ./target --raw --args '@@' &> /dev/null &
done
