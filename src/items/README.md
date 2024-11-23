### Items service

#### Running Locally

```bash
make run
```

#### Running with Docker

1. Build your docker image
```bash
docker build -t items:latest .
```

2. Run your docker container
```bash
docker run -p 8081:8081 -v $(pwd)/db.json:/db.json items:latest
```

3. Send requests to `/items` or `/item?={n}`

```bash
curl http://localhost:8080/api/item\?id\=9 | jq '.'
```

