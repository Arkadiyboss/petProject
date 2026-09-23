FROM golang:1.27.1 AS builder

WORKDIR /app

COPY go.mod go.sum config.json secret.json ./
RUN go mod download
COPY . .

# Ключевое: CGO_ENABLED=0 → статический бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o myapp .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/myapp .
COPY --from=builder /app/config.json .
COPY --from=builder /app/secret.json .

RUN chmod +x /root/myapp

EXPOSE 8080

CMD ["/root/myapp"]