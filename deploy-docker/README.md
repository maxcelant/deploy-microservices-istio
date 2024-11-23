### Deploy Using Docker Compose

1. Make sure you `build` each docker container first. You can do this quickly by running the `./docker-build` script.

```bash
chmod +x scripts/docker-build.sh
./scripts/docker-build
```

1. Run the `compose` command
```bash
docker compose up
```

1. Test it out

```
curl -X POST -H "Content-Type: application/json" -d '{"userId": 1, "itemId": 3}' \
  http://localhost:8082/api/order
```

Expected response:
```json
{
  "id": 1,
  "userId": 1,
  "items": [3],
  "totalPrice": 199.99,
  "status": "OPEN"
}
```

