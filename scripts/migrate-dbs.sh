#!/bin/bash

migrate -database "postgres://users_user:users_pass@localhost:5432/users_db?sslmode=disable" -path ./src/db/migrations/users up
migrate -database "postgres://items_user:items_pass@localhost:5433/items_db?sslmode=disable" -path ./src/db/migrations/items up
migrate -database "postgres://orders_user:orders_pass@localhost:5434/orders_db?sslmode=disable" -path ./src/db/migrations/orders up