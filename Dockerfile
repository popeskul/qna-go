FROM golang:1.18-alpine3.16 AS builder
RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN go get ./...
RUN go get -u github.com/lib/pq
RUN go build -o main cmd/main.go

# Run stage
#FROM alpine:3.16 AS runner
#WORKDIR /app
#COPY --from=builder /app/main .
#COPY .env .
#COPY ./configs ./configs

CMD ["ls", "-l"]
EXPOSE 8080
CMD ["/app/main"]
