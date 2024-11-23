### Using the Items service

1. Build your docker image
```bash
docker build -t items-svc:latest .
```

2. Run your docker container
```bash
docker run -p 8081:8081 -v $(pwd)/db.json:/db.json items-svc:latest
```

3. Send requests to `/items` or `/item?={n}`

```bash
curl http://localhost:8080/api/item\?id\=9 | jq '.'
```

