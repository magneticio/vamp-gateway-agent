#! /bin/bash

# Ensure we have our directories
dir_confd="/etc/confd/conf.d"
dir_templates="/etc/confd/templates"

mkdir -p "$dir_confd"
mkdir -p "$dir_templates"


# Generate config and templates for HAproxy
echo "creating confd configuration and template"
cat <<EOT > "${dir_confd}/haproxy.toml"
[template]
src = "haproxy.tmpl"
dest = "/usr/local/vamp/haproxy.cfg"
keys = [ "${VAMP_KEY_VALUE_STORE_PATH}" ]
check_cmd = "/usr/local/vamp/haproxy-validate.sh {{.src}}"
reload_cmd = "/usr/local/vamp/haproxy-reload.sh"
EOT

cp /usr/local/vamp/haproxy.basic.cfg "${dir_templates}/haproxy.tmpl"
cat <<EOT >> "${dir_templates}/haproxy.tmpl"
{{getv "${VAMP_KEY_VALUE_STORE_PATH}"}}
EOT


# Generate config and template for Filebeat
cat <<EOT > "${dir_confd}/filebeat.toml"
[template]
src = "filebeat.tmpl"
dest = "/usr/local/filebeat/filebeat.yml"
EOT

cat <<EOT >> "${dir_templates}/filebeat.tmpl"
filebeat.prospectors:
- input_type: log
  paths:
    - /var/log/haproxy.log
  json.message_key: message
  json.add_error_key: false
  json.keys_under_root: true

output.elasticsearch:
  hosts: ["$VAMP_ELASTICSEARCH_URL"]
  index: "vamp-vga-$VAMP_NAMESPACE-%{+yyyy-MM-dd}"
  template.path: \${path.config}/filebeat.template.json

path.home: /usr/local/filebeat
path.config: \${path.home}
path.data: \${path.home}/data
path.logs: /var/log

setup.template.name: "vamp-vga-$VAMP_NAMESPACE"
setup.template.pattern: "vamp-vga-$VAMP_NAMESPACE-*"
EOT
