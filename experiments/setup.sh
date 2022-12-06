#!/bin/bash

set -u
## Download Pre-Reqs:
sudo apt-get update && sudo apt-get install -y git make autoconf automake libtool zlib1g-dev wget screen clang clang-tools gcc

## Clone target Repo
git clone --no-checkout https://github.com/glennrp/libpng.git "repo"
git -C "repo" checkout a37d4836519517bdce6cb9d956092321eca3e73b

## Compiler Target
cd repo
CXX="clang++"
CXXFLAGS=" -fsanitize=address"
CXXFLAGS+=" -fno-omit-frame-pointer"
# Edge Coverage
CXXFLAGS+=" -fsanitize-coverage=edge,trace-pc-guard"

autoreconf -f -i
./configure --with-libpng-prefix=MAGMA_ --disable-shared
make -j$(nproc) clean
make -j$(nproc) libpng16.la

cp .libs/libpng16.a "target/"

# build libpng_read_fuzzer.
$CXX $CXXFLAGS -std=c++11 -I. \
     contrib/oss-fuzz/libpng_read_fuzzer.cc \
     -o "target/target" \
     $LDFLAGS .libs/libpng16.a $LIBS -lz

cp target ../
cd .. && rm -rf repo

## Golang:
wget -q https://go.dev/dl/go1.19.3.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz
PATH=$PATH:/usr/local/go/bin

## Compile World:
cd hopper && go build . && mv hopper ../
cd node && go build . && mv node ../../
cd ../../ && rm -rf	hopper/


