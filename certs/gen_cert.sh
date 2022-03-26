#!/bin/sh
openssl req -new -key ./certs/cert.key -subj "/CN=$1" -sha256 | openssl x509 -req -days 3650 -CA ./certs/ca.crt -CAkey certs/ca.key -set_serial "$2" > ./certs/nck.crt