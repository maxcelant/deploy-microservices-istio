FROM golang:1.23.0-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/users/main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /app/main .
COPY --from=builder /app/db.json .

EXPOSE 8080

CMD ["./main"]
