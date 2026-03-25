FROM golang:alpine AS builder
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.20
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY *.sql /migrations/
ENTRYPOINT ["goose", "-dir", "/migrations", "postgres"]
