#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

reset="$(tput sgr0)"
red="$(tput setaf 1)"
green="$(tput setaf 2)"
yellow="$(tput setaf 3)"

# Check dependencies
hash openssl 2> /dev/null || { echo "${red}please install ${green}openssl${reset}"; exit 1; }
hash docker 2> /dev/null || { echo "${red}please install ${green}docker ${yellow}https://www.docker.com${reset}"; exit 1; }

VGA_TAG=${VGA_TAG:-tls}
VGA_BASE_TAG=${VGA_BASE_TAG:-katana}
VGA_DN=${VGA_DN:-localhost}
CLIENT_DN=${CLIENT_DN:-localhost}

TEMP_DIR=${DIR}/.tmp
VGA_CERTS_DIR=${TEMP_DIR}/vga
CLIENT_CERTS_DIR=${TEMP_DIR}/client

set -e

rm -Rf ${TEMP_DIR} && mkdir -p ${VGA_CERTS_DIR} && mkdir -p ${CLIENT_CERTS_DIR}

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
cat > ${VGA_CERTS_DIR}/vga_root.config <<-EOF
[req]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
x509_extensions = v3_ca

[ dn ]
CN = vga_root

[v3_ca]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = vga_root
EOF

cat > ${VGA_CERTS_DIR}/vga.config <<-EOF
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

#Generate server certificate
openssl genrsa -out ${VGA_CERTS_DIR}/vgaRootCA.key 2048
openssl req -x509 -new -nodes -key ${VGA_CERTS_DIR}/vgaRootCA.key -sha256 -days 1024 -config <(cat ${VGA_CERTS_DIR}/vga_root.config) -out ${VGA_CERTS_DIR}/vgaRootCA.crt
openssl genrsa -out ${VGA_CERTS_DIR}/vga.key 2048
openssl req -new -key ${VGA_CERTS_DIR}/vga.key -config <(cat ${VGA_CERTS_DIR}/vga.config) -out ${VGA_CERTS_DIR}/vga.csr
openssl x509 -req -in ${VGA_CERTS_DIR}/vga.csr -CA ${VGA_CERTS_DIR}/vgaRootCA.crt -CAkey ${VGA_CERTS_DIR}/vgaRootCA.key -CAcreateserial -out ${VGA_CERTS_DIR}/vga.crt -days 500 -sha256
cat ${VGA_CERTS_DIR}/vga.crt ${VGA_CERTS_DIR}/vga.key | tee ${VGA_CERTS_DIR}/vgaCertAndKey.crt

step "Creating client certificate"
cat > ${CLIENT_CERTS_DIR}/client_root.config <<-EOF
[req]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
x509_extensions = v3_ca

[ dn ]
CN = client_root

[v3_ca]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = client_root

EOF

cat > ${CLIENT_CERTS_DIR}/client.config <<-EOF
[req]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
x509_extensions = v3_ca

[ dn ]
CN = ${CLIENT_DN}

[v3_ca]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = ${CLIENT_DN}
EOF

openssl genrsa -out ${CLIENT_CERTS_DIR}/clientRootCA.key 2048
openssl req -x509 -new -nodes -key ${CLIENT_CERTS_DIR}/clientRootCA.key -sha256 -days 1024 -config <(cat ${CLIENT_CERTS_DIR}/client_root.config) -out ${CLIENT_CERTS_DIR}/clientRootCA.crt
openssl genrsa -out ${CLIENT_CERTS_DIR}/client.key 2048
openssl req -new -key ${CLIENT_CERTS_DIR}/client.key -config <(cat ${CLIENT_CERTS_DIR}/client.config) -out ${CLIENT_CERTS_DIR}/client.csr
openssl x509 -req -in ${CLIENT_CERTS_DIR}/client.csr -CA ${CLIENT_CERTS_DIR}/clientRootCA.crt -CAkey ${CLIENT_CERTS_DIR}/clientRootCA.key -CAcreateserial -out ${CLIENT_CERTS_DIR}/client.crt -days 500 -sha256

step "Building VGA Docker image"
step "VGA base image   : " ${VGA_BASE_TAG}
step "VGA new image tag: " ${VGA_TAG}
sed -e "s/VGA_TAG/${VGA_BASE_TAG}/g" ${DIR}/Dockerfile > ${TEMP_DIR}/Dockerfile
docker build -t tyga/vamp-gateway-agent:${VGA_TAG} -f ${TEMP_DIR}/Dockerfile ${TEMP_DIR}

finish
