# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./db/migration

RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.20.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/bin/migrate

COPY start.sh .
RUN chmod +x start.sh

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]