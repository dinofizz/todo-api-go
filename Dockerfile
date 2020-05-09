FROM golang:alpine AS build-env
RUN apk update && apk add gcc libc-dev && rm -rf /var/cache/apk/*
WORKDIR /app
ADD . /app
RUN cd /app && go build -o todo-api

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /app/todo-api /app

ENTRYPOINT ["./todo-api"]

