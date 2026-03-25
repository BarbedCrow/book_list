FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server/

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    adduser -D -g '' appuser

COPY --from=builder /server /server

USER appuser
EXPOSE 8080
ENTRYPOINT ["/server"]
