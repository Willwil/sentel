FROM nginx

RUN rm /etc/nginx/conf.d/*

RUN mkdir /etc/nginx/log/

RUN mkdir /etc/nginx/stream_conf.d/
COPY stream_mqtt.conf /etc/nginx/stream_conf.d/
COPY mqtt.js /etc/nginx/stream_conf.d/

COPY nginx.conf /etc/nginx/

EXPOSE 1883
