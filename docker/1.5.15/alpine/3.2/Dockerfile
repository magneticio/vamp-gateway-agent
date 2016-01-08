FROM alpine:3.2

RUN apk --update add musl-dev linux-headers curl gcc pcre-dev make zlib-dev && \
    mkdir /usr/src && \
    curl -fL http://www.haproxy.org/download/1.5/src/haproxy-1.5.15.tar.gz | tar xzf - -C /usr/src && \
    cd /usr/src/haproxy-1.5.15 && \
    make TARGET=linux2628 USE_PCRE=1 USE_ZLIB=1 && \
    make install-bin && \
    cd .. && \
    rm -rf /usr/src/haproxy-1.5.15 && \
    apk del musl-dev linux-headers curl gcc pcre-dev make zlib-dev && \
    apk add musl pcre zlib && \
    rm /var/cache/apk/*
