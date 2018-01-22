FROM golang:alpine
MAINTAINER Kenneth Herner <kherner@navistone.com>

WORKDIR /go/src/github.com/chosenken/prometheus-kairosdb-adapter
COPY . .

RUN go build -o prometheus-kairosdb-adapter

FROM alpine:3.6
MAINTAINER Kenneth Herner <kherner@navistone.com>
COPY --from=0 /go/src/github.com/chosenken/prometheus-kairosdb-adapter/prometheus-kairosdb-adapter /usr/local/bin/prometheus-kairosdb-adapter

EXPOSE 9201
ENTRYPOINT ["prometheus-kairosdb-adapter"]