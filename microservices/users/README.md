### Using the Users service

1. Build your docker image
```bash
docker build -t users-svc:latest .
```

2. Run your docker container
```bash
docker run -p 8080:8080 -v $(pwd)/db.json:/db.json users-svc:latest
```

3. Send requests to `/users` or `/user?={n}`

```bash
curl http://localhost:8080/user\?id\=9 | jq '.'
```

