#!/bin/bash

# List: screen -list
# Detach: Ctrl-a Ctrl-d
# Attach: screen -r <pid>.master
## Spawn Master
export HOPPER_LOG=1
export HOPPER_LOG_INTERVAL=30
export HOPPER_OUT="/proj/hopper-tests-PG0/readelf-dat"

screen -S master -dm ./master -I ./in -H=20
