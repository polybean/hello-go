FROM golang:1.12.7 AS builder
ADD ./src/ /src
WORKDIR /src
RUN go build -v -o hello-server

FROM alpine:3.4
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
EXPOSE 8080
ENV DB db
CMD ["hello-server"]
COPY --from=builder /src/hello-server /usr/local/bin/hello-server
RUN chmod +x /usr/local/bin/hello-server
