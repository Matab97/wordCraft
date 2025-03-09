#!/bin/bash

# Build and deploy with Docker Compose
docker-compose build
docker-compose up -d

echo "Deployment complete!" 