#!/usr/bin/env bash

[ -e /usr/local/vamp/good2go ] || exit 1

haproxy_cfg="/usr/local/vamp/haproxy.cfg"
haproxy_pid="/usr/local/vamp/haproxy.pid"
haproxy_sock="/usr/local/vamp/haproxy.sock"

[ -e ${haproxy_cfg} ] || (echo "haproxy-reload.sh: error: no such file: ${haproxy_cfg}" && exit 1)

[ -e ${haproxy_pid} ] || touch ${haproxy_pid}

if [ -e ${haproxy_sock} ]; then
  haproxy -f ${haproxy_cfg} -p ${haproxy_pid} -x ${haproxy_sock} -sf $(cat ${haproxy_pid})
else
  haproxy -f ${haproxy_cfg} -p ${haproxy_pid} -sf $(cat ${haproxy_pid})
fi
