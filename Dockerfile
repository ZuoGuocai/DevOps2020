FROM golang:1.12.7-alpine3.9
LABEL maintainer="zuoguocai@126.com"  version="2.0" description="ipcat"  

WORKDIR $GOPATH/src/github.com/ipcat
ADD . $GOPATH/src/github.com/ipcat
RUN go build .
EXPOSE 5000
ENTRYPOINT  ["./ipcat"]
