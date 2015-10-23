# Vamp-proxy-agent

Vamp-proxy-agent is a tiny helper agent that provides the following services: 

- read logs from HAproxy over sockets and push them to Logstash over UDP
- read statistics from HAproxy and push them to Logstash over UDP

## Usage

Run `vamp-proxy-agent -h` to display usage instructions:

```bash
$ ./vamp-proxy-agent -h
Usage of ./vamp-proxy-agent:
  -debug
    	Switches on extra log statements
  -haproxyLogSocket string
    	The location of the socket HAproxy logs to (default "/var/run/haproxy.log.sock")
  -haproxyStatsSocket string
    	The location of the HAproxy stats socket (default "/tmp/haproxy.stats.sock")
  -haproxyStatsType string
    	Which stats to read from haproxy: all, frontend, backend or server. (default "all")
  -logHost string
    	Address of the remote Logstash instance (default "127.0.01")
  -logPort int
    	The UDP input port of the remote Logstash instance (default 10002)
  -statsHost string
    	Address of the remote Logstash instance (default "127.0.01")
  -statsPort int
    	The UDP input port of the remote Logstash instance (default 10003)
```

## Example

This example starts a socket at `/var/run/haproxy.sock`. It reads all log lines HAproxy writes to that socket
and passes them to the UDP port 10002 on server 10.4.0.100. Furthermore, it reads the "backend" statistics 
from the `/tmp/haproxy.stats.sock` socket and passes them to port 10003 on the same host.

```bash
$ ./vamp-proxy-agent 
    -logHost 10.4.0.100 \
    -logPort 10002 \
    -haproxyLogSocket /var/run/haproxy.log.sock \
    -statsHost 10.4.0.100 \
    -statsPort 10003 \
    -haproxyStatsSocket /tmp/haproxy.stats.sock \
    -haproxyStatsType backend
   
```

HAproxy should be configured with the following line in its config to actually send the logs to the socket:

```
global
  log /var/run/haproxy.log.sock local2
  stats socket /tmp/haproxy.stats.sock level admin  
```

Logstash should be configure with the following lines in its config to actually read from the UDP sockets
and splits the HAproxy stats in the proper CSV format.

```
input {
  udp {
    port => 10002
    type => haproxy_log
  }
  udp {
    port => 10003
    type => haproxy_stats
  }
}
filter {
  if [type] == "haproxy_stats" {
    csv {
      columns => [pxname,svname,qcur,qmax,scur,smax,slim,stot,bin,bout,dreq,dresp,ereq,econ,eresp,wretr,wredis,status,weight,act,bck,chkfail,chkdown,lastchg,downtime,qlimit,pid,iid,sid,throttle,lbtot,tracked,type,rate,rate_lim,rate_max,check_status,check_code,check_duration,hrsp_1xx,hrsp_2xx,hrsp_3xx,hrsp_4xx,hrsp_5xx,hrsp_other,hanafail,req_rate,req_rate_max,req_tot,cli_abrt,srv_abrt,comp_in,comp_out,comp_byp,comp_rsp,lastsess,last_chk,last_agt,qtime,ctime,rtime,ttime]
    }
  }
}
output {
  stdout { codec => rubydebug }
}
```