# syntax = docker/dockerfile:experimental

# Leave directive above to handle the corner case as listed below

# Client: Docker Engine - Community
#  Version:           19.03.2
#  API version:       1.38 (downgraded from 1.40)
#  Go version:        go1.12.8
#  Git commit:        6a30dfca03
#  Built:             Thu Aug 29 05:26:30 2019
#  OS/Arch:           linux/amd64
#  Experimental:      false

# Server:
#  Engine:
#   Version:          18.06.1-ce
#   API version:      1.38 (minimum version 1.12)
#   Go version:       go1.10.3
#   Git commit:       e68fc7a/18.06.1-ce
#   Built:            Mon Jul  1 18:53:20 2019
#   OS/Arch:          linux/amd64
#   Experimental:     false

# reference: https://dev.to/plutov/docker-and-go-modules-3kkn
FROM golang:1.12.7 AS builder
WORKDIR /src
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
ADD ./ /src
RUN go build -v -o hello-server

FROM alpine:3.4
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
EXPOSE 8080
ENV DB_HOST db
CMD ["hello-server"]
COPY --from=builder /src/hello-server /usr/local/bin/hello-server
RUN chmod +x /usr/local/bin/hello-server
