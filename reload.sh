#!/usr/bin/env bash

configuration=$1
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pid_file=${dir}/haproxy.pid

if [ ! -e ${pid_file} ] ; then
    touch ${pid_file}
fi

haproxy -f ${configuration} -p ${pid_file} -D -st $(cat ${pid_file})
