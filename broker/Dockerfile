FROM alpine

RUN mkdir /sentel
RUN mkdir /sentel/bin
RUN mkdir /sentel/etc
COPY bin/linux_amd64/broker /sentel/bin/
#COPY broker/broker.conf /sentel/etc/

ENTRYPOINT ["/sentel/bin/broker"]
