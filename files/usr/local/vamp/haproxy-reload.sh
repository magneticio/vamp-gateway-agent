#! /bin/bash

handle() { echo "haproxy-reload.sh: got signal"; exit; }
trap handle SIGINT

[[ -e /usr/local/vamp/good2go ]] || exit 1

haproxy_cfg="/usr/local/vamp/haproxy.cfg"
haproxy_pid="/usr/local/vamp/haproxy.pid"

if [[ ! -e $haproxy_cfg ]] ; then
  >&2 echo "haproxy-reload.sh: error: no such file: $haproxy_cfg"
  exit 1
fi

if [[ ! -e $haproxy_pid ]] ; then
  touch $haproxy_pid
fi

# Get all ports we need to add to the INPUT chain
declare -a PORTS=()
regex='^\s*bind 0\.0\.0\.0:([0-9]+)$'
while read line ;do
  if [[ ${line} =~ $regex ]] ;then
    port="${BASH_REMATCH[1]}"
    PORTS+=(${port})
  fi
done < "${haproxy_cfg}"

# for zero downtime HAProxy reload: http://engineeringblog.yelp.com/2015/04/true-zero-downtime-haproxy-reloads.html
# and for this implementation also: https://github.com/mesosphere/marathon-lb/blob/master/service/haproxy/run

for i in "${PORTS[@]}"; do
  iptables -w -I INPUT -p tcp --dport "${i}" --syn -j DROP
done

sleep 0.1

haproxy -f "$haproxy_cfg" -p "$haproxy_pid" -D -st $( cat $haproxy_pid )

for i in "${PORTS[@]}"; do
  iptables -w -D INPUT -p tcp --dport "${i}" --syn -j DROP
done
