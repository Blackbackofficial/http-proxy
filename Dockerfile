FROM golang:1.18-alpine AS lang

ADD . /opt/app
WORKDIR /opt/app

RUN CGO_ENABLED=0 go build ./cmd/main.go

FROM ubuntu:20.04

RUN apt-get -y update
RUN apt-get install -y tzdata

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=lang /opt/app/ .
EXPOSE 8080
ENV PGPASSWORD docker
CMD service postgresql start &&  psql -h localhost -d docker -U docker -p 5432 -a -q -f ./db/init.sql && cd ./certs && chmod 777 ./gen_ca.sh && chmod 777 ./gen_cert.sh && ./gen_ca.sh && cd .. && ./main