FROM redis

RUN apt-get update && apt-get install -y \
  gettext-base \
  && rm -rf /var/lib/apt/lists/*

COPY redis.conf /tmp/

COPY startup.sh /
RUN chmod a+x /startup.sh

CMD ["/startup.sh"]

LABEL MAINTAINER="mihai.oprea@gmail.com"