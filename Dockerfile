FROM golang:1.21-alpine3.17 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY src src
RUN go build -v -o app ./src

FROM alpine:3.17

RUN apk update && apk upgrade

RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]