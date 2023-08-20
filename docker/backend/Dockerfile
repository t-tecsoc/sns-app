# syntax=docker/dockerfile:1
ARG PORT=8080
FROM golang:alpine
ENV APP_DIR=/app

WORKDIR $APP_DIR
COPY ./.env ./

ENV SRC_DIR=backend
WORKDIR $APP_DIR/$SRC_DIR

COPY $SRC_DIR/go.mod ./
COPY $SRC_DIR/go.sum ./
RUN go mod download

COPY $SRC_DIR/./ ./

RUN go build -o /docker-gs-ping

EXPOSE $PORT

CMD [ "/docker-gs-ping" ]