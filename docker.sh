#!/usr/bin/env bash

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

reset=`tput sgr0`
green=`tput setaf 2`
yellow=`tput setaf 3`

version="0.8.0"
target='target'
target_docker=${target}'/docker'
target_go=${target}'/go'
assembly_go='vamp-gateway-agent.tar.gz'

cd ${dir}

function go_build() {
    bin='vamp-gateway-agent'
    export GOOS='linux'
    export GOARCH='amd64'
    echo "${green}building ${GOOS}:${GOARCH} ${yellow}${bin}${reset}"
    rm -rf ${target_go} && mkdir -p ${target_go}
    go build
    mv ${bin} ${target_go} && chmod +x ${target_go}/${bin}
    cp -r ${dir}/configuration/ ${target_go}
    cd ${target_go}
    tar -zcf ${assembly_go} *
    cd ${dir}
}

function docker_rmi {
    echo "${green}removing docker image: $1 ${reset}"
    docker rmi -f $1 2> /dev/null
}

function docker_build {
    echo "${green}appending common code to: $2/Dockerfile ${reset}"
    echo 'RUN mkdir -p /opt/vamp' >> $2/Dockerfile
    echo 'COPY vamp-gateway-agent.tar.gz /opt/vamp/' >> $2/Dockerfile
    echo 'RUN tar -xvzf /opt/vamp/vamp-gateway-agent.tar.gz -C /opt/vamp && rm /opt/vamp/vamp-gateway-agent.tar.gz' >> $2/Dockerfile
    echo 'EXPOSE 1988' >> $2/Dockerfile
    echo 'ENTRYPOINT ["/opt/vamp/vamp-gateway-agent"]' >> $2/Dockerfile

    echo "${green}building docker image: $1 ${reset}"
    docker build -t $1 $2
}

function docker_images {
    arr=$1[@]
    images=("${!arr}")
    pattern=$(printf "\|%s" "${images[@]}")
    pattern=${pattern:2}
    echo "${green}built images:${yellow}"
    docker images | grep 'magneticio/vamp-gateway-agent' | grep ${pattern}
}

echo "${green}cleaning...${reset}"
rm -Rf ${dir}/${target_docker} 2> /dev/null && mkdir -p ${target_docker}

echo "${green}copying files...${reset}"
cp -R ${dir}/docker/* ${dir}/${target_docker}
regex="^${target_docker}\/(.+)\/(.+)\/(.+)\/Dockerfile$"

images=()

for file in `find ${target_docker} | grep Dockerfile`
do
	[[ ${file} =~ $regex ]]
    haproxy_version="${BASH_REMATCH[1]}"
    linux="${BASH_REMATCH[2]}"
    linux_version="${BASH_REMATCH[3]}"
    target=${dir}/${target_docker}/${haproxy_version}/${linux}/${linux_version}
    image=${haproxy_version}-${linux}-${linux_version}
    images+=(${image})
    image_name=magneticio/vamp-gateway-agent_${image}:${version}

    if [[ $* == *--clean* || $* == *-c* ]]
    then
        docker_rmi ${image_name}
    fi
    if [[ $* == *--build* || $* == *-b* ]]
    then
        go_build
        cp -R ${dir}/${target_go}/${assembly_go} ${target} 2> /dev/null
        docker_build ${image_name} ${target}
    fi
done

if [[ $* == *--list* || $* == *-l* ]]
then
    docker_images images
fi

echo "${green}done.${reset}"

