## Vamp Gateway Agent & TLS Termination

### Creating VGA with self-signed TLS certificate

Run `build.sh`, possible parameters (environment variables):

- VGA_TAG=${VGA_TAG:-tls}
- VGA_BASE_TAG=${VGA_BASE_TAG:-katana}
- VGA_DN=${VGA_DN:-localhost}

Script steps:

- creating self-signed certificate
- building the new VGA Docker image with the certificate

Intermediate files including certificate are in `./.tmp` directory.

Another approach is to use official VGA (without building the custom image) and mounting Docker volume with right certificate(s).

### HAProxy configuration

In order to avoid warnings add to global HAProxy configuration:
```
tune.ssl.default-dh-param 2048
```

Assuming `/usr/local/vamp/vga.pem` certificate path, update virtual hosts section: 
```
  bind 0.0.0.0:80
⇒ bind 0.0.0.0:80 ssl crt /usr/local/vamp/vga.pem
```

Also TLS termination can be done differently. For instance just to terminate and proxy to a gateway port, replace virtual hosts part with:
```
### BEGIN - TLS TERMINATION

frontend tls_termination

  bind 0.0.0.0:443 ssl crt /usr/local/vamp/vga.pem
  mode http

  option httplog
  log-format """{\"ci\":\"%ci\",\"cp\":%cp,\"t\":\"%t\",\"ft\":\"%ft\",\"b\":\"%b\",\"s\":\"%s\",\"Tq\":%Tq,\"Tw\":%Tw,\"Tc\":%Tc,\"Tr\":%Tr,\"Tt\":%Tt,\"ST\":%ST,\"B\":%B,\"CC\":\"%CC\",\"CS\":\"%CS\",\"tsc\":\"%tsc\",\"ac\":%ac,\"fc\":%fc,\"bc\":%bc,\"sc\":%sc,\"rc\":%rc,\"sq\":%sq,\"bq\":%bq,\"hr\":\"%hr\",\"hs\":\"%hs\",\"r\":%{+Q}r}"""
  
  use_backend tls_termination

backend tls_termination

  balance roundrobin
  mode http

  option forwardfor
  http-request set-header X-Forwarded-Port %[dst_port]
  
  # server: sava/80
  server tls_termination 127.0.0.1:80
  
### END - TLS TERMINATION
```

Note `127.0.0.1:80` where `80` is a Vamp gateway port.

### Example

Let's assume the following:

- VGA is accessible on `aaa-bbb-ccc.eu-west-1.elb.amazonaws.com`.
- we have `sava` deployed (port 9050) and virtual hosts enabled
- it should work: `curl -H 'Host: 9050.sava.vamp' http://aaa-bbb-ccc.eu-west-1.elb.amazonaws.com`

Now, setup the TLS using self-signed certificate:

- run: `export VGA_DN=*.elb.amazonaws.com && ./build.sh`
- we just created (by default): `magneticio/vamp-gateway-agent:tls`, redeploy VGA using that image
- it should work as before: `curl -H 'Host: 9050.sava.vamp' http://aaa-bbb-ccc.eu-west-1.elb.amazonaws.com`
- go to Vamp VGA template and update `bind 0.0.0.0:80` to `bind 0.0.0.0:443 ssl crt /usr/local/vamp/vga.pem`
- check now (notice `https`): `curl --insecure -H 'Host: 9050.sava.vamp' https://aaa-bbb-ccc.eu-west-1.elb.amazonaws.com`
