#!/bin/bash

echo "Building microservices..."
docker build -t items:latest ./src/items
echo "[1/3] Items Service Completed"
docker build -t orders:latest ./src/orders
echo "[2/3] Orders Service Completed"
docker build -t users:latest ./src/users
echo "[3/3] Users Service Completed"
echo "Done."