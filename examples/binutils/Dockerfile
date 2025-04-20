FROM --platform=linux/amd64 hopper-node:latest

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl python3 tar xz-utils texinfo zlib1g-dev build-essential file
# Build binutils
WORKDIR /binutils
COPY . .

RUN ./build-binutils.py -t arm-linux-gnueabi

RUN cp ./install/bin/arm-linux-gnueabi-readelf /readelf_target

RUN rm -rf /hopper/examples/binutils/binutils-* /hopper/examples/binutils/install /hopper/examples/binutils/build
