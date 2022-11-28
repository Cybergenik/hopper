#!/bin/bash

for ((i=1;i<=$1;i++))
do
  go run -race node/node.go -I $i -T test/target --args "@@" &
done
