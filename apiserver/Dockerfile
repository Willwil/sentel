FROM alpine
RUN mkdir /sentel
RUN mkdir /sentel/bin
RUN mkdir /sentel/etc
COPY bin/linux_amd64/apiserver /sentel/bin/
#COPY apiserver/apiserver.conf /sentel/etc/
COPY apiserver/swagger.yaml /sentel/etc/

ENTRYPOINT ["/sentel/bin/apiserver"]
