#!/usr/bin/env bash

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

reset=`tput sgr0`
green=`tput setaf 2`
yellow=`tput setaf 3`

target='target'
target_docker=${target}'/docker'
project="vamp-gateway-agent"

if [ "$(git describe --tags)" = "$(git describe --abbrev=0 --tags)" ]; then
    version="$( git describe --tags )"
    docker_image_name="magneticio/${project}:${version}"
else
    version="katana [$( git describe --tags )]"
    docker_image_name="magneticio/${project}:katana"
fi

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
        -r|--remove)
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

function build_help() {
    echo "${green}Usage of $0:${reset}"
    echo "${yellow}  -h|--help   ${green}Help.${reset}"
    echo "${yellow}  -l|--list   ${green}List built Docker images.${reset}"
    echo "${yellow}  -r|--remove ${green}Remove Docker image.${reset}"
    echo "${yellow}  -m|--make   ${green}Make Docker image files.${reset}"
    echo "${yellow}  -b|--build  ${green}Build Docker image.${reset}"
}

function docker_make {
    echo ${version} > ${dir}/${target_docker}/version
    cp ${dir}/Dockerfile ${dir}/${target_docker}/Dockerfile
    cp -Rf ${dir}/files ${dir}/${target_docker}
}

function docker_build {
    echo "${green}building docker image: $1 ${reset}"
    docker build -t $1 $2
}

function docker_rmi {
    echo "${green}removing docker image: $1 ${reset}"
    docker rmi -f $1 2> /dev/null
}

function docker_image {
    echo "${green}built images:${yellow}"
    docker images | grep "magneticio/${project}"
}

function process() {

    echo "${green}version: ${version}${reset}"

    rm -Rf ${dir}/${target} 2> /dev/null && mkdir -p ${dir}/${target_docker}

    if [ ${flag_make} -eq 1 ]; then
        docker_make
    fi

    if [ ${flag_clean} -eq 1 ]; then
        docker_rmi ${docker_image_name}
    fi

    if [ ${flag_build} -eq 1 ]; then
        cd ${dir}/${target_docker}
        docker_build ${docker_image_name} .
    fi

    if [ ${flag_list} -eq 1 ]; then
        docker_image
    fi

    echo "${green}done.${reset}"
}

parse_command_line $@

if [ ${flag_help} -eq 1 ] || [[ $# -eq 0 ]]; then
    build_help
fi

if [ ${flag_list} -eq 1 ] || [ ${flag_clean} -eq 1 ] || [ ${flag_make} -eq 1 ] || [ ${flag_build} -eq 1 ]; then
    process
fi
