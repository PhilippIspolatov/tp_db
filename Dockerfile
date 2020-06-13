FROM golang:1.13-stretch AS builder

# Building project
WORKDIR /build

COPY . .
RUN go build -v ./cmd/app/main.go

FROM ubuntu

# Expose server & database ports
EXPOSE 5000
EXPOSE 5432

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get -y update && apt-get install -y --no-install-recommends apt-utils postgresql-12;

USER postgres

ENV PGPASSWORD="techdb"


# Create & configure database
COPY /sql/init.sql .
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql -f ./init.sql -d docker &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/12/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/12/main/postgresql.conf

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

COPY ./sql/init.sql ./assets/db/postgres/base.sql

# Copying built binary
COPY --from=builder /build/main .
CMD service postgresql start && ./main