FROM golang:1.18-alpine

ENV GO111MODULE=on

ENV GOPATH=/

COPY ./ ./

EXPOSE 8080

RUN apk --no-cache add curl openssl

RUN chmod 777 ./certs/gen_ca.sh && chmod 777 ./certs/gen_cert.sh

RUN ./certs/gen_ca.sh

RUN go mod download

RUN go mod tidy

RUN go build -o http-proxy ./cmd/main.go

CMD ./http-proxy
