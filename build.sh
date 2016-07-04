#!/usr/bin/env bash

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

reset=`tput sgr0`
green=`tput setaf 2`
yellow=`tput setaf 3`

version="$( git describe --tags )"
target='target'
target_vamp=${target}'/vamp'
target_docker=${target}'/docker'
project="vamp-gateway-agent"
docker_image_name="magneticio/${project}:${version}"

cd ${dir}

function parse_command_line() {
    flag_help=0
    flag_list=0
    flag_clean=0
    flag_make=0
    flag_build=0
    flag_build_all=0

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
        -a|--all)
        flag_build_all=1
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
    echo "${yellow}  -m|--make   ${green}Build the binary and copy it to the Docker directories.${reset}"
    echo "${yellow}  -b|--build  ${green}Build Docker image.${reset}"
    echo "${yellow}  -a|--all    ${green}Build all binaries, by default only linux:amd64.${reset}"
}

function go_make() {
    cd ${dir}
    rm -Rf ${dir}/${target_vamp}

    echo "${green}executing ${yellow}godep restore${reset}"
    go get github.com/tools/godep
    godep restore
    go install

    for goos in darwin linux windows; do
      for goarch in 386 amd64; do

        if [ ${flag_build_all} -eq 1 ] || [[ ${goos} == "linux" && ${goarch} == "amd64" ]]; then

          cd ${dir}
          mkdir ${dir}/${target_vamp}

          export GOOS=${goos}
          export GOARCH=${goarch}

          echo "${green}building ${yellow}${project}_${version}_${goos}_${goarch}${reset}"

          CGO_ENABLED=0 go build -ldflags "-X main.version=${version}" -a -installsuffix cgo

          if [ "${goos}" == "windows" ]; then
              mv ${dir}/${project}.exe ${target_vamp}
          else
              mv ${dir}/${project} ${target_vamp} && chmod +x ${target_vamp}/${project}
          fi

          assembly_go="${project}_${version}_${goos}_${goarch}.tar.gz"

          cp -f ${dir}/reload.sh ${dir}/validate.sh ${dir}/haproxy.basic.cfg ${dir}/${target_vamp}
          cd ${dir}/${target} && tar -zcf ${assembly_go} vamp
          mv ${dir}/${target}/${assembly_go} ${dir}/${target_docker} 2> /dev/null

          rm -Rf ${dir}/${target_vamp} 2> /dev/null

        fi
      done
    done
}

function docker_make {

    append_to=${dir}/${target_docker}/Dockerfile
    cat ${dir}/Dockerfile | grep -v ADD | grep -v ENTRYPOINT > ${append_to}

    echo "${green}appending common code to: ${append_to} ${reset}"
    function append() {
        printf "\n$1\n" >> ${append_to}
    }

    append "ADD ${project}_${version}_linux_amd64.tar.gz /usr/local"
    append "ENTRYPOINT [\"/usr/local/vamp/${project}\"]"
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

    rm -Rf ${dir}/${target} 2> /dev/null && mkdir -p ${dir}/${target_docker} && mkdir -p ${target_vamp}

    if [ ${flag_make} -eq 1 ]; then
        docker_make
        go_make
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
