#!/usr/bin/env bash

configuration=$1
pid_file=/tmp/haproxy.pid

if [ ! -e ${pid_file} ] ; then
    touch ${pid_file}
fi

PORTS=()

regex='^\s*bind 0\.0\.0\.0:([0-9]+)$'
while read line
do
    if [[ ${line} =~ $regex ]]
    then
        port="${BASH_REMATCH[1]}"
        PORTS+=(${port})
    fi
done < "${configuration}"

# for zero downtime HAProxy reload: http://engineeringblog.yelp.com/2015/04/true-zero-downtime-haproxy-reloads.html
# and for this implementation also: https://github.com/mesosphere/marathon-lb/blob/master/service/haproxy/run

for i in "${PORTS[@]}"; do
  iptables -w -I INPUT -p tcp --dport ${i} --syn -j DROP
done

sleep 0.1

haproxy -f ${configuration} -p ${pid_file} -D -st $(cat ${pid_file})

for i in "${PORTS[@]}"; do
  iptables -w -D INPUT -p tcp --dport ${i} --syn -j DROP
done

sleep 1
