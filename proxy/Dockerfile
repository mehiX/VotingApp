FROM nginx:latest

RUN rm -f /etc/nginx/conf.d/default.conf

COPY index.html /usr/share/nginx/html/
COPY app.template /etc/nginx/conf.d/

CMD envsubst < /etc/nginx/conf.d/app.template > /etc/nginx/conf.d/app.conf && exec nginx -g 'daemon off;'

LABEL MAINTAINER="mihai.oprea@gmail.com"