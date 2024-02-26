#!/bin/bash

# Stop and delete containers
docker-compose down
docker stop inmemory

# Delete all our images
docker rmi shortlink-grpc
docker rmi postgres:15
docker rmi inmemory