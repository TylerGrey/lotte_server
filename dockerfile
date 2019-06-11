FROM golang:1.12.5 AS build-env
ENV GO111MODULE=on

ENV APP_DIR $GOPATH/src/github.org/TylerGrey/lotte_server

ENTRYPOINT ["/api"]
ADD . $APP_DIR

RUN go mod download
RUN cd $APP_DIR && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api

RUN ls -al /

# Multi stage
FROM alpine:3.7 as prod_img
COPY --from=0 /api /api
RUN ls -al /

ENTRYPOINT ["/api"]