#!/usr/bin/env bash

configuration=$1
pid_file=/usr/local/vamp/haproxy.pid

if [ ! -e ${pid_file} ] ; then
    touch ${pid_file}
fi

haproxy -f ${configuration} -p ${pid_file} -D -st $(cat ${pid_file})
