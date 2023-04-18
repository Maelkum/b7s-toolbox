# Readme

This docker compose setup can be used to start multiple b7s nodes on the same machine.
The `docker-compose.yaml` defines four nodes - one head node and three worker nodes.
All nodes use predefined keys so they'll have the same identities each time they are started.

## Dockerfile

The docker compose file uses a Docker image `b7s-node` for individual nodes/services.
The accompanying Dockerfile can be used to build the image.
