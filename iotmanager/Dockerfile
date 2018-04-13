FROM alpine
RUN mkdir /sentel
COPY bin/linux_amd64/iotmanager /sentel/
# COPY iotmanager/iotmanager.conf /etc/sentel

ENTRYPOINT ["/sentel/iotmanager"]
