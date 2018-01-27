FROM alpine
RUN mkdir /sentel
COPY bin/linux_amd64/iothub /sentel/
# COPY broker/iothub.conf /etc/sentel

ENTRYPOINT ["/sentel/iothub"]
