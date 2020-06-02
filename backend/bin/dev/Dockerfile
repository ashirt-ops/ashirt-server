FROM golang:1.14.4-alpine

RUN apk add --no-cache git inotify-tools g++ && \
    rm -rf /var/cache/apk/* && \
    go get -v github.com/rubenv/sql-migrate/...

RUN mkdir /app
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

CMD ["backend/run-dev.sh"]
