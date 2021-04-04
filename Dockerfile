# Basic dockerfile
FROM alpine:latest

# Set up golang
ENV GOPATH /golang
RUN apk update && apk add go
RUN mkdir /golang

# Set up (compile + install) binaries
RUN mkdir /app
COPY . /app
RUN cd /app; go install ./cmd/todo

# Command
CMD [ "/golang/bin/todo" ]
