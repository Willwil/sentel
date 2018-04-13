FROM alpine
RUN mkdir /sentel
# RUN mkdir /etc/sentel
COPY bin/linux_amd64/apiserver /sentel/
# COPY apiserver/apiserver.conf /etc/sentel

ENTRYPOINT ["/sentel/apiserver"]
