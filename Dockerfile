FROM alpine:3.4

ENV VAMP_GATEWAY_VERSION=0.9.0

RUN set -ex && \
    apk --update add musl-dev linux-headers curl gcc pcre-dev make zlib-dev && \
    mkdir /usr/src && \
    curl -fL http://www.haproxy.org/download/1.6/src/haproxy-1.6.6.tar.gz | tar xzf - -C /usr/src && \
    cd /usr/src/haproxy-1.6.6 && \
    make TARGET=linux2628 USE_PCRE=1 USE_ZLIB=1 && \
    make install-bin && \
    cd .. && \
    rm -rf /usr/src/haproxy-1.6.6 && \
    apk del musl-dev linux-headers curl gcc pcre-dev make zlib-dev && \
    apk add musl pcre zlib && \
    rm /var/cache/apk/*

EXPOSE 1988

ADD https://bintray.com/artifact/download/magnetic-io/downloads/vamp-gateway-agent/vamp-gateway-agent_${VAMP_GATEWAY_VERSION}_linux_amd64.tar.gz /usr/local

ENTRYPOINT ["/usr/local/vamp/vamp-gateway-agent"]
