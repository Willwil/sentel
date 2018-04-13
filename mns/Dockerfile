FROM alpine
RUN mkdir /sentel
# RUN mkdir /etc/sentel
COPY bin/linux_amd64/mns /sentel/
# COPY mns/mns.conf /etc/sentel

ENTRYPOINT ["/sentel/mns"]
