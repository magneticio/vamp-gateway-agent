# Vamp Gateway Agent & TLS Termination

## Creating VGA with a TLS certificate

In this example self-signed certificate will be used.
Run `build.sh`, possible parameters (environment variables):

- VGA_TAG=${VGA_TAG:-tls}
- VGA_BASE_TAG=${VGA_BASE_TAG:-katana}
- VGA_DN=${VGA_DN:-localhost}

## Updating Vamp HAProxy template
