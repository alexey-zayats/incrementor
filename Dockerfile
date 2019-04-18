FROM golang:alpine as builder

RUN apk -U --no-cache add git make protobuf protobuf-dev \
    && go get -u github.com/golang/protobuf/protoc-gen-go

ENV GOROOT /usr/local/go

ADD . /src/incrementor
WORKDIR /src/incrementor

RUN make

FROM alpine

COPY --from=builder /src/incrementor/bin/incrementor /app/incrementor

WORKDIR /app

VOLUME ["/app/config"]

ENTRYPOINT ["/app/incrementor", "server"]
