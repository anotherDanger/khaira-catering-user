FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

RUN mkdir -p /tmp/uploads && chmod 775 /tmp/uploads

COPY --from=builder /app/app .

EXPOSE 8083

ENTRYPOINT ["./app"]