FROM --platform=linux/amd64 ubuntu:latest 

WORKDIR /app

COPY ./raft-node /app/node

# run script
ENTRYPOINT [ "/app/node" ]
