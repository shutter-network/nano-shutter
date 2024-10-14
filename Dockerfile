FROM golang:1.23-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git build-base
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o nano-shutter

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app .

EXPOSE 5000

CMD ["./nano-shutter"]