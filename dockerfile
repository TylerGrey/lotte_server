FROM golang:1.12.5 AS build-env
ENV GO111MODULE=on

ENV APP_DIR $GOPATH/src/github.org/TylerGrey/lotte_server

ENTRYPOINT ["/api"]
ADD . $APP_DIR
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

RUN go mod download
RUN cd $APP_DIR && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api


RUN ls -al /

# Multi stage
FROM alpine:3.7 as prod_img
COPY --from=0 /api /api
COPY --from=0 /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
RUN ls -al /

ENTRYPOINT ["/api"]