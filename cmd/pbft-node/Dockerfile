FROM --platform=linux/amd64 ubuntu:latest 

WORKDIR /app

## curl, unzip other utilities
RUN apt-get update && \
  apt-get install --no-install-recommends --assume-yes curl unzip pv ca-certificates gnupg2 net-tools jq

  
# get the runtime
# RUN curl -o ./runtime.tar.gz -sSL https://github.com/blocklessnetwork/runtime/releases/download/v0.0.13/blockless-runtime.ubuntu-20.04.x86_64.tar.gz
# RUN mkdir /app/runtime && tar -xvkf ./runtime.tar.gz -C /app/runtime

# libssl 1.1
RUN curl -o ./libssl.deb -sSL http://archive.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.1_1.1.0g-2ubuntu4_amd64.deb
RUN dpkg -i ./libssl.deb

COPY ./pbft-node /app/node

# run script
ENTRYPOINT [ "/app/node" ]
