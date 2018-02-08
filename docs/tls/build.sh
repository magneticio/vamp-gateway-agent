#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

reset="$(tput sgr0)"
red="$(tput setaf 1)"
green="$(tput setaf 2)"
yellow="$(tput setaf 3)"

# Check dependencies
hash openssl 2> /dev/null || { echo "${red}please install ${green}openssl${reset}"; exit 1; }
hash docker 2> /dev/null || { echo "${red}please install ${green}docker ${yellow}https://www.docker.com${reset}"; exit 1; }

TEMP_DIR=${DIR}/.tmp
VGA_TAG=${VGA_TAG:-tls}
VGA_BASE_TAG=${VGA_BASE_TAG:-katana}
VGA_DN=${VGA_DN:-localhost}

set -e

rm -Rf ${TEMP_DIR} && mkdir -p ${TEMP_DIR}

function step {
  echo "${green}$1${yellow}$2"
}

function finish {
  step "Done."
  exit 0
}

function panic {
  echo "${red}ERROR: $1${reset}"
  exit 1
}

function terminated {
  printf ${reset}
}

trap terminated EXIT

step "Creating VGA certificate"
cat > ${TEMP_DIR}/vga.csr <<-EOF
[req]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
x509_extensions = v3_ca

[ dn ]
CN = ${VGA_DN}

[v3_ca]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = ${VGA_DN}
EOF
openssl req -x509 -nodes -newkey rsa:2048 -keyout ${TEMP_DIR}/vga_key.pem -out ${TEMP_DIR}/vga_cert.pem -days 3650 -config <(cat ${TEMP_DIR}/vga.csr)
cat ${TEMP_DIR}/vga_cert.pem ${TEMP_DIR}/vga_key.pem | tee ${TEMP_DIR}/vga.pem

step "Building VGA Docker image"
step "VGA base image   : " ${VGA_BASE_TAG}
step "VGA new image tag: " ${VGA_TAG}
sed -e "s/VGA_TAG/${VGA_BASE_TAG}/g" ${DIR}/Dockerfile > ${TEMP_DIR}/Dockerfile
docker build -t magneticio/vamp-gateway-agent:${VGA_TAG} -f ${TEMP_DIR}/Dockerfile ${TEMP_DIR}

finish
