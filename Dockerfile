FROM alpine:3.5

ADD version vamp-gateway-agent.sh reload.sh validate.sh haproxy.basic.cfg /usr/local/vamp/
ADD https://github.com/kelseyhightower/confd/releases/download/v0.11.0/confd-0.11.0-linux-amd64 /usr/bin/confd

RUN set -ex && \
    apk --update add bash iptables musl-dev linux-headers curl gcc pcre-dev make zlib-dev dnsmasq && \
    mkdir /usr/src && \
    curl -fL http://www.haproxy.org/download/1.7/src/haproxy-1.7.1.tar.gz | tar xzf - -C /usr/src && \
    cd /usr/src/haproxy-1.7.1 && \
    make TARGET=linux2628 USE_PCRE=1 USE_ZLIB=1 && \
    make install-bin && \
    cd .. && \
    rm -rf /usr/src/haproxy-1.7.1 && \
    apk del musl-dev linux-headers curl gcc pcre-dev make zlib-dev && \
    apk add musl pcre zlib && \
    rm /var/cache/apk/* && \
    echo "port=5353" > /etc/dnsmasq.conf && \
    chmod u+x /usr/bin/confd && \
    chmod u+x /usr/local/vamp/vamp-gateway-agent.sh /usr/local/vamp/reload.sh /usr/local/vamp/validate.sh

EXPOSE 1988

ENTRYPOINT ["/usr/local/vamp/vamp-gateway-agent.sh"]

