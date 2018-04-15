FROM alpine
RUN mkdir /sentel
RUN mkdir /sentel/bin
RUN mkdir /sentel/etc
COPY bin/linux_amd64/whaler /sentel/bin/
#COPY whaler/whaler.conf /sentel/etc/

ENTRYPOINT ["/sentel/bin/whaler"]
