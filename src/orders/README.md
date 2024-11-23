### Orders service

#### Running Locally

```bash
make run
```

#### Running with Docker
1. Build your docker image
```bash
docker build -t orders:latest .
```

2. Run your docker container
```bash
docker run -p 8082:8082 
```

3. Send requests to `/users` or `/user?={n}`

```bash
curl http://localhost:8082/api/order\?id\=2 | jq '.'
```

