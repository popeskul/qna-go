FROM golang:1.18-alpine3.16 AS builder

WORKDIR /app
COPY . .

RUN GOOS=linux go build -o ./bin/app ./cmd/main.go

# Run stage
FROM alpine:latest AS runner

COPY --from=builder /app/bin/app/ .
COPY --from=builder /app/.env .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./app"]