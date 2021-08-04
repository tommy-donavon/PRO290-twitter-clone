FROM ubuntu:trusty

VOLUME /etc/consul-template/templates
VOLUME /etc/consul-template/conf

COPY conf/* /etc/consul-template/conf/
COPY templates/* /etc/consul-template/templates/

ADD nginx.conf /usr/nginx/nginx.conf

CMD tail -f /dev/null