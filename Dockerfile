FROM alpine:3.6

# https://github.com/peterbourgon/runsvinit
ENV RUNSVINIT_URL=https://github.com/peterbourgon/runsvinit/releases/download/v2.0.0/runsvinit-linux-amd64.tgz

ENV HAPROXY_VER=1.7.5
ENV HAPROXY_URL=http://www.haproxy.org/download/1.7/src/haproxy-${HAPROXY_VER}.tar.gz

ENV CONFD_URL=https://github.com/kelseyhightower/confd/releases/download/v0.13.0/confd-0.13.0-linux-amd64

ENV FILEBEAT_VER=5.1.2
ENV FILEBEAT_URL=https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-${FILEBEAT_VER}-linux-x86_64.tar.gz

RUN set -xe \
    && apk add --no-cache \
      bash \
      curl \
      iptables \
      musl \
      pcre \
      rsyslog \
      runit \
      zlib \
      jq \
    && curl --location --silent --show-error $RUNSVINIT_URL --output - | tar zxf - -C /sbin \
    && chown 0:0 /sbin/runsvinit \
    && chmod 0775 /sbin/runsvinit \
    \
    && curl --location --silent --show-error --output /usr/bin/confd $CONFD_URL \
    && chmod 0755 /usr/bin/confd \
    \
    && apk --no-cache --virtual=build-deps add \
      gcc \
      linux-headers \
      make \
      musl-dev \
      pcre-dev \
      zlib-dev \
    && mkdir /usr/src \
    && curl -fL $HAPROXY_URL | tar xzf - -C /usr/src \
    && cd /usr/src/haproxy-${HAPROXY_VER} \
    && make TARGET=linux2628 USE_PCRE=1 USE_ZLIB=1 \
    && make install-bin \
    && cd .. \
    && rm -rf /usr/src/haproxy-${HAPROXY_VER} \
    && apk del build-deps \
    \
    && curl --location --silent --show-error $FILEBEAT_URL --output - | tar zxf - -C /tmp \
    && mv /tmp/filebeat-${FILEBEAT_VER}-linux-x86_64/filebeat /usr/local/bin/ \
    && rm -rf /tmp/filebeat-${FILEBEAT_VER}-linux-x86_64

RUN ALPINE_GLIBC_BASE_URL="https://github.com/sgerrand/alpine-pkg-glibc/releases/download" && \
    ALPINE_GLIBC_PACKAGE_VERSION="2.23-r3" && \
    ALPINE_GLIBC_BASE_PACKAGE_FILENAME="glibc-$ALPINE_GLIBC_PACKAGE_VERSION.apk" && \
    ALPINE_GLIBC_BIN_PACKAGE_FILENAME="glibc-bin-$ALPINE_GLIBC_PACKAGE_VERSION.apk" && \
    ALPINE_GLIBC_I18N_PACKAGE_FILENAME="glibc-i18n-$ALPINE_GLIBC_PACKAGE_VERSION.apk" && \
    apk add --no-cache --virtual=.build-dependencies wget ca-certificates && \
    wget \
        "https://raw.githubusercontent.com/andyshinn/alpine-pkg-glibc/master/sgerrand.rsa.pub" \
        -O "/etc/apk/keys/sgerrand.rsa.pub" && \
    wget \
        "$ALPINE_GLIBC_BASE_URL/$ALPINE_GLIBC_PACKAGE_VERSION/$ALPINE_GLIBC_BASE_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_BASE_URL/$ALPINE_GLIBC_PACKAGE_VERSION/$ALPINE_GLIBC_BIN_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_BASE_URL/$ALPINE_GLIBC_PACKAGE_VERSION/$ALPINE_GLIBC_I18N_PACKAGE_FILENAME" && \
    apk add --no-cache \
        "$ALPINE_GLIBC_BASE_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_BIN_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_I18N_PACKAGE_FILENAME" && \
    \
    rm "/etc/apk/keys/sgerrand.rsa.pub" && \
    /usr/glibc-compat/bin/localedef --force --inputfile POSIX --charmap UTF-8 C.UTF-8 || true && \
    echo "export LANG=C.UTF-8" > /etc/profile.d/locale.sh && \
    \
    apk del glibc-i18n && \
    \
    rm "/root/.wget-hsts" && \
    apk del .build-dependencies && \
    rm \
        "$ALPINE_GLIBC_BASE_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_BIN_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_I18N_PACKAGE_FILENAME"

ENV LANG=C.UTF-8


EXPOSE 1988

ADD files/ /
ADD version /usr/local/vamp/version

RUN chmod +x /usr/local/vamp/tokenrenewer.sh

ENTRYPOINT ["/sbin/runsvinit"]
