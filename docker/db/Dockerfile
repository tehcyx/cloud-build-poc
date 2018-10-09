FROM postgres:9.4-alpine

ARG buildDateLabel
LABEL buildDateLabel=$buildDateLabel

RUN mkdir -p /docker-entrypoint-initdb.d
COPY init-user-db.sh /docker-entrypoint-initdb.d