FROM alpine
RUN mkdir /sentel
RUN mkdir /sentel/bin
RUN mkdir /sentel/etc
COPY bin/linux_amd64/mns /sentel/bin/
#COPY mns/mns.conf /sentel/etc/

ENTRYPOINT ["/sentel/bin/mns"]
