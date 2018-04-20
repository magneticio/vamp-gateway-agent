# Vamp Gateway Agent

[![Join the chat at https://gitter.im/magneticio/vamp](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/magneticio/vamp?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Docker](https://img.shields.io/badge/docker-images-blue.svg)](https://hub.docker.com/r/magneticio/vamp-gateway-agent/tags/)

Based on Vamp gateways, Vamp generates HAProxy configuration and stores it to KV store.

Vamp Gateway Agent:
- reads the [HAProxy](http://www.haproxy.org/) configuration using [confd](https://github.com/kelseyhightower/confd)
- appends it to the base configuration `haproxy.basic.cnf`
- if new configuration is valid, VGA reloads HAProxy with as little client traffic interruption as possible

In addition to this VGA also:
- send HAProxy log to Elasticsearch using [Filebeat](https://www.elastic.co/products/beats/filebeat)
- handle and recover from ZooKeeper, etcd, Consul and Vault outages without interrupting the haproxy process and client requests
- does Vault token renewal if needed

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

## Building Docker images

`make` targets:
- `version` - displaying version (tag)
- `clean` - removing temporal build directory `./target`
- `purge` - running `clean` and removing image `magneticio/vamp-gateway-agent:${version}`
- `build` - copying files to `./target` directory and building the image `magneticio/vamp-gateway-agent:${version}`
- `default` - `clean build`


## Additional documentation and examples

- [TLS](https://github.com/magneticio/vamp-gateway-agent/tree/master/docs/tls)
