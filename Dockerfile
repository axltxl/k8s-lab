# FIXME: doc me
FROM alpine:latest

ENV GOPATH /golang

RUN apk update && apk add go
RUN mkdir /app; mkdir /golang

COPY . /app
RUN cd /app; go install ./cmd/todo

CMD [ "/golang/bin/todo" ]
