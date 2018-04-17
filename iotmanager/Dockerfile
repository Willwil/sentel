FROM ubuntu
RUN mkdir /sentel
RUN mkdir /sentel/bin
RUN mkdir /sentel/etc
COPY bin/linux_amd64/iotmanager /sentel/bin/
COPY deploy/wait-for-it.sh /sentel/bin/
#COPY iotmanager/iotmanager.conf /sentel/etc/

CMD ["/sentel/bin/iotmanager"]
