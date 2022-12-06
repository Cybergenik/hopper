#!/bin/bash

for ((i=$1;i<=$2;i++))
do
	./node -I $i -T ./target/target -M <change_me> --stdin &
done
