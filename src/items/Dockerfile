FROM golang:1.23.0-bullseye AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/items/main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]
