FROM centos:7

RUN yum install -y epel-release yum install -y gcc make pcre-devel zlib-devel

RUN curl -fL http://www.haproxy.org/download/1.5/src/haproxy-1.5.14.tar.gz | tar xzf - -C /usr/src && \
cd /usr/src/haproxy-1.5.14 && \
make TARGET=linux2628 USE_PCRE=1 USE_ZLIB=1 && \
make install && \
cd .. && \
rm -rf /usr/src/haproxy-1.5.14

