FROM golang:1.12 AS builder
WORKDIR /go/src/app

COPY Makefile .
COPY go.mod .
COPY go.sum .
RUN make download

COPY . .
RUN make build && cp bin/prometheus-rough-proxy /usr/bin/prometheus-rough-proxy

###############################################

FROM ubuntu:18.04
COPY --from=builder /usr/bin/prometheus-rough-proxy /usr/bin/prometheus-rough-proxy

ENTRYPOINT ["/usr/bin/prometheus-rough-proxy"]

