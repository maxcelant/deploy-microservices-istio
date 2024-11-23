#!/bin/bash

echo "Building microservices..."
docker build -t items:latest ./src/shopclub/v1/items
echo "[1/3] Items Service Completed"
docker build -t orders:latest ./src/shopclub/v1/orders
echo "[2/3] Orders Service Completed"
docker build -t users:latest ./src/shopclub/v1/users
echo "[3/3] Users Service Completed"
echo "Done."