#!/usr/bin/env bash

haproxy_cfg="/usr/local/vamp/haproxy.cfg"

[[ -n $1 ]] && haproxy_cfg="$1"

haproxy -c -f "$haproxy_cfg"
