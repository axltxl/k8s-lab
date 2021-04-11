# Basic dockerfile
FROM golang:1.16-alpine

# Set up golang
ENV GOPATH /golang
RUN mkdir /golang

# Set up (compile + install) binaries
RUN mkdir /app
COPY . /app
RUN cd /app; go install ./src/cmd/todod

# Command
CMD [ "/golang/bin/todod" ]
