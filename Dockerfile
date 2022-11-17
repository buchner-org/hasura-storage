FROM golang:1.19 AS builder

WORKDIR /build

COPY . .

ENV VERSION=0.1
ENV MODULE=github.com/nhost/hasura-storage
RUN apt-get update \
  && apt-get install -y ca-certificates libvips-dev \
  && go mod download \
  && go build -ldflags "-X ${MODULE}/controller.buildVersion=${VERSION}" -o /build/hasura-storage-bin \
  && strip hasura-storage-bin

FROM debian:latest

WORKDIR /app

ENV TMPDIR=/
ENV MALLOC_ARENA_MAX=2

COPY --from=builder /build/hasura-storage-bin /app/hasura-storage

RUN apt-get update \
  && apt-get install -y ca-certificates libvips \
  && useradd -d /app -s /bin/bash storage \
  && touch config-file-dev.yml

USER storage

CMD [ "/app/hasura-storage-bin", "serve" ]

