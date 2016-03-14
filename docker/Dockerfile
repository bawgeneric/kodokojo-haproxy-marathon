FROM haproxy:1.6

ADD kodokojo-haproxy-marathon /usr/bin/kodokojo-haproxy-marathon
ADD start.sh /start.sh
ADD haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg

RUN chmod +x /usr/bin/kodokojo-haproxy-marathon
RUN chmod +x /start.sh
RUN mkdir -p /usr/local/etc/haproxy/ssl/

EXPOSE 80 443 4444

ENTRYPOINT /start.sh