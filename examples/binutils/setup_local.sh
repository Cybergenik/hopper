#!/bin/bash

set -e

# Update and install git and curl
sudo apt-get update && sudo apt-get install -y git curl clang clang-tools gcc wget

cd ~
# Download Go
wget -q https://go.dev/dl/go1.19.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Clone and build Hopper
git clone --branch stress_tests --single-branch https://github.com/Cybergenik/hopper.git
cd hopper/

go build .
cd node/
go build .

# Copy run scripts
cd ~
cp hopper/hopper ./master
cp hopper/node/node .
cp ~/hopper/examples/binutils/readelf/dist/master_local.sh .
cp ~/hopper/examples/binutils/readelf/dist/node_local.sh .
rm -rf hopper
