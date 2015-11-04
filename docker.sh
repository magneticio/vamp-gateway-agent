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

function parse_command_line() {
    flag_help=0
    flag_list=0
    flag_clean=0
    flag_make=0
    flag_build=0

    for key in "$@"
    do
    case ${key} in
        -h|--help)
        flag_help=1
        ;;
        -l|--list)
        flag_list=1
        ;;
        -c|--clean)
        flag_clean=1
        ;;
        -m|--make)
        flag_make=1
        ;;
        -b|--build)
        flag_make=1
        flag_build=1
        ;;
        *)
        ;;
    esac
    done
}

function help() {
    echo "${green}Usage of $0:${reset}"
    echo "${yellow}  -h|--help   ${green}Help.${reset}"
    echo "${yellow}  -l|--list   ${green}List all available images.${reset}"
    echo "${yellow}  -c|--clean  ${green}Remove all available images.${reset}"
    echo "${yellow}  -m|--make   ${green}Build vamp-gateway-agent binary.${reset}"
    echo "${yellow}  -b|--build  ${green}Build all available images.${reset}"
}

function go_build() {
    bin='vamp-gateway-agent'
    export GOOS='linux'
    export GOARCH='amd64'
    echo "${green}building ${GOOS}:${GOARCH} ${yellow}${bin}${reset}"
    rm -rf ${target_go} && mkdir -p ${target_go}

    go get github.com/tools/godep
    godep restore
    go install
    go build

    mv ${bin} ${target_go} && chmod +x ${target_go}/${bin} && cp ${dir}/haproxy.cfg ${target_go}/.
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

function process() {
    rm -Rf ${dir}/${target_docker} 2> /dev/null && mkdir -p ${target_docker}
    cp -R ${dir}/docker/* ${dir}/${target_docker}

    if [ ${flag_make} -eq 1 ]; then
        go_build
    fi

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

        if [ ${flag_clean} -eq 1 ]; then
            docker_rmi ${image_name}
        fi
        if [ ${flag_build} -eq 1 ]; then
            cp -R ${dir}/${target_go}/${assembly_go} ${target} 2> /dev/null
            docker_build ${image_name} ${target}
        fi
    done

    if [ ${flag_list} -eq 1 ]; then
        docker_images images
    fi

    echo "${green}done.${reset}"
}

parse_command_line $@

if [ ${flag_help} -eq 1 ] || [[ $# -eq 0 ]]; then
    help
fi

if [ ${flag_list} -eq 1 ] || [ ${flag_clean} -eq 1 ] || [ ${flag_build} -eq 1 ]; then
    process
fi
