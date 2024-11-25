### Debugging Postgres Containers

Make one database for each microservice.
```bash
# For the users service
docker run --name postgres-users -e POSTGRES_USER=users_user -e POSTGRES_PASSWORD=users_pass -e POSTGRES_DB=users_db -p 5432:5432 -d postgres

# For the items service
docker run --name postgres-items -e POSTGRES_USER=items_user -e POSTGRES_PASSWORD=items_pass -e POSTGRES_DB=items_db -p 5433:5432 -d postgres

# For the orders service
docker run --name postgres-orders -e POSTGRES_USER=orders_user -e POSTGRES_PASSWORD=orders_pass -e POSTGRES_DB=orders_db -p 5434:5432 -d postgres
```

Run the migration for the database

```bash
migrate -database "postgres://users_user:users_pass@localhost:5432/users_db?sslmode=disable" -path ./src/db/migrations/users up
```

Check that it was successful

```bash
docker exec -it <postgres-container> psql -U <posgres-user> -d <database-name>
```

Delete all records in a database table

```bash
docker exec -it cb2907531b5a psql -U users_user -d users_db
psql (17.2 (Debian 17.2-1.pgdg120+1))
Type "help" for help.

users_db=# DELETE FROM users;
DELETE 3
users_db=# \q
```

