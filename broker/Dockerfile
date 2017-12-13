FROM alpine
RUN mkdir /sentel
# RUN mkdir /etc/sentel
COPY bin/linux_amd64/broker /sentel/
# COPY broker/broker.conf /etc/sentel

ENTRYPOINT ["/sentel/broker", "-t", "suny", "-p", "product1", "-P", "tcp"]
