FROM ubuntu
RUN mkdir /sentel
RUN mkdir /etc/sentel
COPY bin/linux_amd64/sentel-broker /sentel
COPY broker/broker.conf /etc/sentel

ENTRYPOINT ["/sentel/sentel-broker"]
