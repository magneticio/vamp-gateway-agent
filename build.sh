#!/usr/bin/env bash

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

reset=`tput sgr0`
green=`tput setaf 2`
yellow=`tput setaf 3`

version="0.8.5"
target='target'
target_vamp=${target}'/vamp'
target_haproxy=${target}'/haproxy'
assembly_go='vamp.tar.gz'

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

function build_help() {
    echo "${green}Usage of $0:${reset}"
    echo "${yellow}  -h|--help   ${green}Help.${reset}"
    echo "${yellow}  -l|--list   ${green}List all available images.${reset}"
    echo "${yellow}  -c|--clean  ${green}Remove all available images.${reset}"
    echo "${yellow}  -m|--make   ${green}Build vamp-gateway-agent binary and copy it to Docker directories.${reset}"
    echo "${yellow}  -b|--build  ${green}Build all available images.${reset}"
}

function go_build() {
    cd ${dir}
    bin='vamp-gateway-agent'
    export GOOS='linux'
    export GOARCH='amd64'
    echo "${green}building ${GOOS}:${GOARCH} ${yellow}${bin}${reset}"
    rm -rf ${target_vamp} && mkdir -p ${target_vamp}

    go get github.com/tools/godep
    godep restore
    go install
    CGO_ENABLED=0 go build -a -installsuffix cgo

    mv ${bin} ${target_vamp} && chmod +x ${target_vamp}/${bin}
}

function docker_rmi {
    echo "${green}removing docker image: $1 ${reset}"
    docker rmi -f $1 2> /dev/null
}

function docker_make {

    append_to=$2/Dockerfile
    mkdir -p "$(dirname "${append_to}")" && touch "${append_to}"
    cat $1 | grep -v ADD | grep -v ENTRYPOINT > ${append_to}
    echo "${green}appending common code to: ${append_to} ${reset}"
    function append() {
        printf "\n$1\n" >> ${append_to}
    }

    append "ADD ${assembly_go} /opt"
    append "ENTRYPOINT [\"/opt/vamp/vamp-gateway-agent\"]"
}

function docker_build {
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
    rm -Rf ${dir}/${target} 2> /dev/null && mkdir -p ${dir}/${target_haproxy} && cd ${dir}/haproxy

    if [ ${flag_make} -eq 1 ]; then
        go_build
    fi

    cd ${dir}
    regex="^${dir}\/haproxy\/(.+)\/Dockerfile$"
    images=()

    for file in `find ${dir}/haproxy | grep Dockerfile`
    do
      [[ ${file} =~ $regex ]]
        haproxy_version="${BASH_REMATCH[1]}"
        docker_path=${dir}/${target_haproxy}/${haproxy_version}
        image=${haproxy_version}
        images+=(${image})
        image_name=magneticio/vamp-gateway-agent_${image}:${version}

        if [ ${flag_make} -eq 1 ]; then
          docker_make ${file} ${docker_path}
          cp -f ${dir}/reload.sh ${dir}/${target_vamp}
          cp -f ${dir}/validate.sh ${dir}/${target_vamp}
          cp -f ${dir}/haproxy/${haproxy_version}/haproxy.basic.cfg ${dir}/${target_vamp}
          cd ${dir}/${target} && tar -zcf ${assembly_go} vamp
          mv ${dir}/${target}/${assembly_go} ${docker_path} 2> /dev/null
        fi

        if [ ${flag_clean} -eq 1 ]; then
            docker_rmi ${image_name}
        fi

        if [ ${flag_build} -eq 1 ]; then
            docker_build ${image_name} ${docker_path}
        fi
    done

    if [ ${flag_list} -eq 1 ]; then
        docker_images images
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
