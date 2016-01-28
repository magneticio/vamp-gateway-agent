FROM centos:7

RUN set -ex && \
    yum -y update && \
    yum install -y epel-release gcc make pcre-devel zlib-devel && \
    curl -fL http://www.haproxy.org/download/1.6/src/haproxy-1.6.3.tar.gz | tar xzf - -C /usr/src && \
    cd /usr/src/haproxy-1.6.3 && \
    make TARGET=linux2628 USE_PCRE=1 USE_ZLIB=1 && \
    make install && \
    cd .. && \
    rm -rf /usr/src/haproxy-1.6.3 && \
    yum remove -y epel-release gcc make && \
    yum clean all
