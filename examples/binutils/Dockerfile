FROM --platform=linux/amd64 hopper-node:latest
RUN apt-get install -y curl python3 tar xz-utils texinfo zlib1g-dev build-essential file
# Build binutils
RUN cd /hopper/examples/binutils/ && ./build-binutils.py -t arm-linux-gnueabi
RUN cd /hopper && cp /hopper/examples/binutils/install/bin/arm-linux-gnueabi-readelf ./readelf_target
RUN rm -rf /hopper/examples/binutils/binutils-* /hopper/examples/binutils/install /hopper/examples/binutils/build
ENV PATH=$PATH:/usr/local/go/bin

