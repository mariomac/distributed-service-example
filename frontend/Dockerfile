FROM golang:1.20 AS builder

WORKDIR /src

COPY . .

RUN make compile

FROM ubuntu:latest

WORKDIR /
COPY --from=builder /src/bin/frontend /frontend

ENTRYPOINT ["/frontend"]
