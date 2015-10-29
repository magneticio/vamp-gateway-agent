FROM ubuntu:14.04

RUN apt-get update && apt-get install -y curl gcc make libpcre3-dev zlib1g-dev

RUN curl -fL http://www.haproxy.org/download/1.5/src/haproxy-1.5.14.tar.gz | tar xzf - -C /usr/src && \
cd /usr/src/haproxy-1.5.14 && \
make TARGET=linux26 USE_PCRE=1 USE_ZLIB=1 USE_LINUX_SPLICE=1 USE_LINUX_TPROXY=1 && \
make install-bin && \
cd .. && \
rm -rf /usr/src/haproxy-1.5.14

