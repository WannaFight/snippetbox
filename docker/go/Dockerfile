FROM golang:1.22.1

WORKDIR /app

COPY . /app

RUN mkdir tls && cd tls && go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

RUN go install github.com/cosmtrek/air@latest
