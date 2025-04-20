 #!/bin/bash
set -e

# Master config
HOPPER_OUT="/hopper_out"
CORPUS_PATH="/corpus"
HAVOC=2

# create docker subnet
docker network create hopper-jq-subnet &> /dev/null

docker run --rm -it \
    --name hopper-master-jq \
    --env TERM \
    --env HOPPER_OUT=$HOPPER_OUT \
    --env HOPPER_LOG=1 \
    --env HOPPER_LOG_INTERVAL=10 \
    --volume $(pwd)$HOPPER_OUT:$HOPPER_OUT \
    --volume $(pwd)$CORPUS_PATH:$CORPUS_PATH \
    --network hopper-jq-subnet \
    --publish 6969:6969 \
    hopper-jq:latest \
    bash -c "hopper-master -I ${CORPUS_PATH} -H ${HAVOC} -P 6969"

## Clean up subnet
docker network rm -f hopper-jq-subnet
