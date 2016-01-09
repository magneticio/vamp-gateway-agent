FROM alpine:3.3

RUN apk --update add musl-dev linux-headers curl gcc pcre-dev make zlib-dev && \
    mkdir /usr/src && \
    curl -fL http://www.haproxy.org/download/1.6/src/haproxy-1.6.3.tar.gz | tar xzf - -C /usr/src && \
    cd /usr/src/haproxy-1.6.3 && \
    make TARGET=linux2628 USE_PCRE=1 USE_ZLIB=1 && \
    make install-bin && \
    cd .. && \
    rm -rf /usr/src/haproxy-1.6.3 && \
    apk del musl-dev linux-headers curl gcc pcre-dev make zlib-dev && \
    apk add musl pcre zlib && \
    rm /var/cache/apk/*
