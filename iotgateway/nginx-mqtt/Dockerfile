FROM centos

# install software
COPY ./nginx-mqtt/nginx.repo /etc/yum.repos.d/
RUN yum -y install supervisor \
    && yum -y install nginx \
    && yum -y install nginx-module-njs \
    && yum clean all

# config nginx mqtt
COPY nginx-mqtt/etc /

# config sentel-gateway
RUN mkdir /sentel
COPY etc/template.conf /sentel/
COPY bin/linux_amd64/sentel-gateway /sentel/


EXPOSE 80 1883

CMD ["/usr/bin/supervisord"]
