FROM alpine:latest
MAINTAINER jenson.zuo@qq.com
WORKDIR $GOPATH/src/github.com/cloustone/sentel/bin
ADD /sentel/sentel-broker  $GOPATH/src/github.com/cloustone/sentel/bin/sentel-broker
add /etc/sentel/ $GOPATH/src/github.com/cloustone/etc/sentel/broker.conf
EXPOSE 8080

ENTRYPOINT ["/sentel/sentel-broker"]

