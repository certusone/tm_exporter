FROM golang:stretch as builder

COPY . /opt
WORKDIR /opt

RUN make build

FROM scratch

COPY --from=builder /opt/bin/tm_exporter /

CMD ["/tm_exporter"]
