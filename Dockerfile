FROM --platform=linux/amd64 debian:stable
COPY . /hopper
EXPOSE 6969/tcp
## Install clang utils
RUN apt-get update && apt-get install -y \
    clang \
    clang-tools \
    gcc \
    wget
## Install Go
RUN wget -q https://go.dev/dl/go1.19.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
RUN cd /hopper && go build ./cmd/hopper-master && go build ./cmd/hopper-node

