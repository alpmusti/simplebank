# Build stage
FROM golang:1.17-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz migrate

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY db/migration ./migration
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh", "/app/main" ]