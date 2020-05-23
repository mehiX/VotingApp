FROM nginx:latest

RUN rm -f /etc/nginx/conf.d/default.conf

COPY index.html /usr/share/nginx/html/
COPY *.conf /etc/nginx/conf.d/

LABEL MAINTAINER="mihai.oprea@gmail.com"