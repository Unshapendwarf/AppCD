FROM golang:1.13.7 AS builder

MAINTAINER Gyutae Park (gtp7473@sk.com)

RUN mkdir -p /build
WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod vendor
COPY . .

RUN go build -o cmd

RUN mkdir -p /dist
WORKDIR /dist
RUN cp /build/cmd ./appcd

RUN ldd appcd | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/



FROM alpine:3.11.3

RUN mkdir -p /app
WORKDIR /app

COPY --chown=0:0 --from=builder /dist /app/
RUN mv lib/x86_64-linux-gnu/* /lib/. && cp -r ./lib64 /

EXPOSE 8080

ENTRYPOINT ["/app/appcd"]
