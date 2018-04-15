FROM alpine
RUN mkdir /sentel
# RUN mkdir /etc/sentel
COPY bin/linux_amd64/whaler /sentel/
# COPY whaler/whaler.conf /etc/sentel

ENTRYPOINT ["/sentel/whaler"]
