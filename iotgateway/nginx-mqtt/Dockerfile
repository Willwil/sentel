FROM centos

# install software
COPY iotgateway/nginx-mqtt/nginx.repo /etc/yum.repos.d/
RUN yum -y install python-setuptools \
    && easy_install supervisor \
    && yum -y install nginx \
    && yum -y install nginx-module-njs \
    && yum clean all

# config nginx mqtt
COPY iotgateway/nginx-mqtt/etc /

# config sentel-gateway
RUN mkdir /sentel
COPY iotgateway/etc/template.conf /sentel/
COPY bin/linux_amd64/iotgateway /sentel/


EXPOSE 80 1883

CMD ["/usr/bin/supervisord"]
