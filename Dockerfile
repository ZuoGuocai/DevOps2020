
FROM golang:latest
LABEL maintainer="zuoguocai@126.com"  version="2.0" description="ipcat"


ENV ELASTIC_APM_SERVICE_NAME=ipcat
ENV ELASTIC_APM_SERVER_URL=https://8b06fec588334601ba91e8ad7fe235c3.apm.eastus2.azure.elastic-cloud.com:443
ENV ELASTIC_APM_SECRET_TOKEN=JOFwFHBYdXzAbIMUYP

WORKDIR /myapp
ADD ./ipcat /myapp/ipcat
EXPOSE 5000
ENTRYPOINT  ["./ipcat"]







#FROM golang:1.12.7-alpine3.9
#LABEL maintainer="zuoguocai@126.com"  version="2.0" description="ipcat"  

#WORKDIR $GOPATH/src/github.com/ipcat
#ADD . $GOPATH/src/github.com/ipcat
#RUN go build .
#EXPOSE 5000
#ENTRYPOINT  ["./ipcat"]
