FROM alpine
RUN mkdir /sentel
RUN mkdir /sentel/bin
RUN mkdir /sentel/etc
COPY bin/linux_amd64/iotmanager /sentel/bin/
#COPY iotmanager/iotmanager.conf /sentel/etc/

ENTRYPOINT ["/sentel/bin/iotmanager"]
