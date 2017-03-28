# Vamp Gateway Agent

[![Join the chat at https://gitter.im/magneticio/vamp](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/magneticio/vamp?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Docker](https://img.shields.io/badge/docker-images-blue.svg)](https://hub.docker.com/r/magneticio/vamp-gateway-agent/tags/)
[![Download](https://api.bintray.com/packages/magnetic-io/downloads/vamp-gateway-agent/images/download.svg) ](https://bintray.com/magnetic-io/downloads/vamp-gateway-agent/_latestVersion)

[HAProxy](http://www.haproxy.org/) is a tcp/http load balancer, the purpose of this agent is to: 

- read the HAProxy configuration using [confd](https://github.com/kelseyhightower/confd) and reload HAProxy on each configuration change with as little client traffic interruption as possible.
- send HAProxy log to Logstash/Elasticsearch.
- handle and recover from ZooKeeper, etcd, Consul and Logstash outages without interrupting the haproxy process and client requests.

Vamp generated HAProxy configuration will be appended to base configuration `haproxy.basic.cnf`.
This allows using different base configuration if needed.

## Usage

Following environment variables are mandatory:

- `VAMP_KEY_VALUE_STORE_TYPE <=> confd -backend`
- `VAMP_KEY_VALUE_STORE_CONNECTION <=> confd -node`
- `VAMP_KEY_VALUE_STORE_PATH <=> key used by confd`
- `VAMP_ELASTICSEARCH_URL <=> http://elasticsearch:9200`

Example:

```
docker run -e VAMP_KEY_VALUE_STORE_TYPE=zookeeper \
           -e VAMP_KEY_VALUE_STORE_CONNECTION=localhost:2181 \
           -e VAMP_KEY_VALUE_STORE_PATH=/vamp/gateways/haproxy/1.6 \
           -e VAMP_ELASTICSEARCH_URL=http://localhost:9200 \
           magneticio/vamp-gateway-agent:katana
```

Available Docker images can be found at [Docker Hub](https://hub.docker.com/r/magneticio/vamp-gateway-agent/).

### Domain name resolver

To enable dnsmasq to resolve virtual hosts, pass the following environment variables to the Docker container:

- `VAMP_VGA_DNS_ENABLE` Set to non-empty value to enable 
- `VAMP_VGA_DNS_PORT` Listening port, default: 5353

### Metrics

The Vamp gateway agent docker image uses Metricbeat to collect performance metrics and ship them off to Elasticsearch. 
By default the [system module](https://www.elastic.co/guide/en/beats/metricbeat/current/metricbeat-module-system.html) is configured to store metrics, with the additional tags to ease filtering:

- `vamp`
- `gateway`

 
## Building Docker Images

```shell
make
```

Docker images after the build: `magneticio/vamp-gateway-agent:katana`

For more details on available targets see the contents of the `Makefile`.
