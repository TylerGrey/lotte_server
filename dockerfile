FROM golang:1.12.5
ENV GO111MODULE=on

ENV APP_DIR $GOPATH/src/github.org/TylerGrey/lotte_server
COPY .env .
RUN export $(cat .env)

ENTRYPOINT ["/api"]
ADD . $APP_DIR

RUN go mod download
RUN cd $APP_DIR && go build -o /api
