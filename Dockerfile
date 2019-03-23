FROM golang:1.11-alpine

WORKDIR /go/src/github.com/isogram/clean-golang

COPY . /go/src/github.com/isogram/clean-golang

ENV BUILD_DEP "git curl zip make build-base"
ENV GO111MODULE=on  

RUN apk add --no-cache rsyslog supervisor py-pip tzdata && \
    apk add --no-cache ${BUILD_DEP} && \
    make all && \
    cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    echo "Asia/Jakarta" >  /etc/timezone && \
    apk del ${BUILD_DEP}

COPY _build/supervisord.conf /etc/supervisor.d/user-service.ini

# RUN touch /tmp/supervisor.sock
# RUN chmod 777 /tmp/supervisor.sock

EXPOSE 3000

CMD ["/usr/bin/supervisord", "-nc", "/etc/supervisord.conf"]
