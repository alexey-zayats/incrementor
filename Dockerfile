FROM golang:alpine as builder

RUN apk -U --no-cache add git make protobuf protobuf-dev \
    && go get -u github.com/golang/protobuf/protoc-gen-go \
    && go get -u github.com/hashicorp/go-multierror \
    && go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq \
    && go build -tags 'postgres' -o /go/bin/migrate github.com/golang-migrate/migrate/cli

ENV GOROOT /usr/local/go

ADD . /src/incrementor
WORKDIR /src/incrementor

RUN make

FROM alpine

COPY --from=builder /src/incrementor/bin/incrementor /app/incrementor
COPY --from=builder /go/bin/migrate /bin/migrate
COPY --from=builder /src/incrementor/scripts/migrate.sh /app/migrate
COPY --from=builder /src/incrementor/migrations /app/migrations

RUN apk -U --no-cache add bash \
    && chmod +x /app/migrate \
    && chmod +x /bin/migrate

WORKDIR /app

VOLUME ["/app/config"]

ENTRYPOINT ["/app/incrementor", "server"]
