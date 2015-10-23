# Vamp-proxy-agent

Vamp-proxy-agent is a tiny helper agent that provides the following services: 

- read logs from HAproxy over sockets and push them to Logstash
- read statistics from HAproxy and push them to Logstash

## Usage

Run `vamp-proxy-agent -h` to display usage instructions:

```bash
$ ./vamp-proxy-agent -h
Usage of ./vamp-proxy-agent:
  -debug
    	Switches on extra log statements
  -haproxyLogSocket string
    	The file location of the socket HAproxy logs to (default "/var/run/haproxy.log.sock")
  -logstashHost string
    	Address of the remote Logstash instance (default "127.0.01")
  -logstashPort int
    	The UDP input port of the remote Logstash instance (default 10002)
```

## Example

This example starts a socket at /var/run/haproxy.sock. It reads all log lines HAproxy writes to that socket
and passes them to the UDP port 10003 on server 10.4.0.100.

```bash
$ ./vamp-proxy-agent -logstashHost 10.4.0.100 -logstashPort 10003 -haproxyLogSocket /var/run/haproxy.log.sock
```

HAproxy should be configured with the following line in its config to actually send the logs to the socket:

```
global
  log /var/run/haproxy.log.sock local2
```

Logstash should be configure with the following line in its config to actually read from the UDP socket:

```
input { 
  udp {
    port => 10003
  }
}
```