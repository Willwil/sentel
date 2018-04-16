FROM ubuntu:14.04
RUN mkdir /sentel
RUN mkdir /sentel/bin
RUN mkdir /sentel/etc
COPY bin/linux_amd64/whaler /sentel/bin/
COPY scripts/wait-for-it.sh /sentel/bin/
#COPY whaler/whaler.conf /sentel/etc/

CMD ["/sentel/bin/whaler"]
