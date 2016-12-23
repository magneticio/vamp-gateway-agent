#!/usr/bin/env sh

echo "
██╗   ██╗ █████╗ ███╗   ███╗██████╗      ██████╗  █████╗ ████████╗███████╗██╗    ██╗ █████╗ ██╗   ██╗
██║   ██║██╔══██╗████╗ ████║██╔══██╗    ██╔════╝ ██╔══██╗╚══██╔══╝██╔════╝██║    ██║██╔══██╗╚██╗ ██╔╝
██║   ██║███████║██╔████╔██║██████╔╝    ██║  ███╗███████║   ██║   █████╗  ██║ █╗ ██║███████║ ╚████╔╝
╚██╗ ██╔╝██╔══██║██║╚██╔╝██║██╔═══╝     ██║   ██║██╔══██║   ██║   ██╔══╝  ██║███╗██║██╔══██║  ╚██╔╝
 ╚████╔╝ ██║  ██║██║ ╚═╝ ██║██║         ╚██████╔╝██║  ██║   ██║   ███████╗╚███╔███╔╝██║  ██║   ██║
  ╚═══╝  ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝          ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝ ╚══╝╚══╝ ╚═╝  ╚═╝   ╚═╝
"

: "${VAMP_KEY_VALUE_STORE_PATH?not provided.}"
: "${VAMP_KEY_VALUE_STORE_TYPE?not provided.}"
: "${VAMP_KEY_VALUE_STORE_CONNECTION?not provided.}"

printf "VERSION: " && cat /usr/local/vamp/version
/usr/bin/confd -version

echo "VAMP_KEY_VALUE_STORE_TYPE       : ${VAMP_KEY_VALUE_STORE_TYPE}"
echo "VAMP_KEY_VALUE_STORE_CONNECTION : ${VAMP_KEY_VALUE_STORE_CONNECTION}"
echo "VAMP_KEY_VALUE_STORE_PATH       : ${VAMP_KEY_VALUE_STORE_PATH}"

mkdir -p /usr/local/vamp/confd/conf.d
mkdir -p /usr/local/vamp/confd/templates

echo "creating confd configuration and template"
cat <<EOT > /usr/local/vamp/confd/conf.d/workflow.toml
[template]
src = "haproxy.tmpl"
dest = "/usr/local/vamp/haproxy.cfg"
keys = [ "${VAMP_KEY_VALUE_STORE_PATH}" ]
check_cmd = "/usr/local/vamp/validate.sh {{.src}}"
reload_cmd = "/usr/local/vamp/reload.sh /usr/local/vamp/haproxy.cfg"
EOT
cp /usr/local/vamp/haproxy.basic.cfg /usr/local/vamp/confd/templates/haproxy.tmpl
cat <<EOT >> /usr/local/vamp/confd/templates/haproxy.tmpl
{{getv "${VAMP_KEY_VALUE_STORE_PATH}"}}
EOT

echo "running confd to retrieve HAProxy configuration"
/usr/bin/confd -interval 5 \
        -backend ${VAMP_KEY_VALUE_STORE_TYPE} \
        -node ${VAMP_KEY_VALUE_STORE_CONNECTION} \
        -confdir /usr/local/vamp/confd
