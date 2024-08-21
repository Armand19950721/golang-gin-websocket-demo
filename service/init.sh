#!/bin/bash
# Build the Docker image
docker build -f Dockerfile.build -t myapp-build .

# Create a temporary container to copy the executable
docker create --name temp-container myapp-build

# Copy the executable from the temporary container to the host machine
docker cp temp-container:/main ./main

# Remove the temporary container
docker rm temp-container
