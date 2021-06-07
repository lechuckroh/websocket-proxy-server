# build
FROM golang:1.16 as build
LABEL image=websocket-proxy-server-builder
WORKDIR /usr/src/myapp
COPY . /usr/src/myapp
RUN go mod vendor && \
    CGO_ENABLED=0 GO111MODULE=on GOARCH=amd64 GOOS=linux go build -mod=vendor -v -o ws-proxy

# release
FROM alpine:latest as release
LABEL image=websocket-proxy-server
LABEL MAINTAINER="Lechuck Roh <lechuckroh@gmail.com>"
RUN mkdir -p /app/scripts
COPY --from=build /usr/src/myapp/ws-proxy /app

WORKDIR /app
VOLUME ["/app/scripts"]

CMD /app/ws-proxy
