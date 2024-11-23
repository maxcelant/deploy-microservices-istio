
```bash
postgres=# CREATE DATABASE sampledb;
CREATE DATABASE
postgres=# CREATE USER sampledb_admin WITH PASSWORD 'password';
CREATE ROLE
postgres=# GRANT ALL PRIVILEGES ON DATABASE sampledb TO sampledb_admin;
GRANT
```