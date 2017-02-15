#! /bin/bash

# Ensure we have our directories
dir_confd="/usr/local/vamp/confd/conf.d"
dir_templates="/usr/local/vamp/confd/templates"

mkdir -p "$dir_confd"
mkdir -p "$dir_templates"


# Generate config and templates for HAproxy
echo "vamp-gateway-agent/run: creating confd configuration and template"
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
  index: "vamp-vga-%{+yyyy-MM-dd}"
  template.path: \${path.config}/filebeat.template.json

path.home: /usr/local/filebeat
path.config: \${path.home}
path.data: \${path.home}/data
path.logs: /var/log
EOT


# Generate config and templates for Metricbeat
cat <<EOT > "${dir_confd}/metricbeat.toml"
[template]
src = "metricbeat.tmpl"
dest = "/usr/local/metricbeat/metricbeat.yml"
EOT

cat <<EOT >> "${dir_templates}/metricbeat.tmpl"
metricbeat.modules:
- module: system
  metricsets:
    - cpu         # CPU stats
    - load        # System Load stats
    - filesystem  # Per filesystem stats
    - fsstat      # File system summary stats
    - memory      # Memory stats
    - network     # Network stats
    - process     # Per process stats
  enabled: true
  period: 10s
  processes: ['.*']

output.elasticsearch:
  hosts: ["$VAMP_ELASTICSEARCH_URL"]
  index: "vamp-vwa-%{+yyyy-MM-dd}"
  template.path: /usr/local/metricbeat/metricbeat.template.json

path.home: /usr/local/metricbeat
path.config: \${path.home}
path.data: \${path.home}/data
path.logs: /var/log
EOT

