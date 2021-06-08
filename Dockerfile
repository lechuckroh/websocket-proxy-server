# build
FROM golang:1.16 as build
LABEL image=websocket-proxy-server-builder
WORKDIR /usr/src/myapp
COPY . /usr/src/myapp
RUN go mod vendor && \
    go build -ldflags "-s -w" -v -o ws-proxy

# release
FROM golang:1.16 as release
LABEL image=websocket-proxy-server
LABEL MAINTAINER="Lechuck Roh <lechuckroh@gmail.com>"
RUN mkdir -p /app/scripts
COPY --from=build /usr/src/myapp/ws-proxy /app

WORKDIR /app
VOLUME ["/app/scripts"]

CMD /app/ws-proxy
