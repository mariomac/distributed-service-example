FROM golang:1.20 AS builder

WORKDIR /src

COPY . .

RUN make compile

FROM busybox:1.35-glibc

WORKDIR /
COPY --from=builder /src/bin/backend /backend

ENTRYPOINT ["/backend"]
