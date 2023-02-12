#!/bin/bash

set -e

# Update and install git and curl
sudo apt-get update && sudo apt-get install git curl

cd ~
# Install docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Clone and build Hopper
git clone https://github.com/Cybergenik/hopper.git
cd hopper/
sudo docker build -t hopper-node .
cd examples/binutils/
sudo docker build -t hopper-readelf .
cd ~
