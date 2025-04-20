FROM hopper-node:latest

RUN apt-get update && apt-get install -y --no-install-recommends \
    clang-14 clang-tools-14 libclang-14-dev llvm-14-dev lld-14 libclang-rt-14-dev \
    && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/stedolan/jq.git /jq

WORKDIR /jq

RUN git submodule update --init --recursive

RUN cd vendor/oniguruma && ./autogen.sh && cd ../..

RUN autoreconf -fi
RUN CC=clang \
    CXX=clang++ \
    CFLAGS="-fsanitize=address -fno-omit-frame-pointer -fsanitize-address-use-after-scope -fsanitize-coverage=edge,trace-pc-guard" \
    LDFLAGS="-fsanitize=address" \
    ./configure

RUN make -j$(nproc)

RUN cp /jq/jq /usr/local/bin/
