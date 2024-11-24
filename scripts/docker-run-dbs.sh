#!/bin/bash

echo "Starting postgres containers for each service"
docker run --name postgres-users -e POSTGRES_USER=users_user -e POSTGRES_PASSWORD=users_pass -e POSTGRES_DB=users_db -p 5432:5432 -d postgres
echo "[1/3] User database created"
docker run --name postgres-items -e POSTGRES_USER=items_user -e POSTGRES_PASSWORD=items_pass -e POSTGRES_DB=items_db -p 5433:5432 -d postgres
echo "[2/3] Items database created"
docker run --name postgres-orders -e POSTGRES_USER=orders_user -e POSTGRES_PASSWORD=orders_pass -e POSTGRES_DB=orders_db -p 5434:5432 -d postgres
echo "[3/3] Orders database created"
echo "Done."