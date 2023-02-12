# Header Dump

![Yea it's a stock photo...](/.assets/icon.png)

Does what it says on the tin. Developped as a debugging tool, this traefik middleware sits in the middle of your route and mercilessly dumps all request headers into traefik logs in a date-time-log format. The `Prefix` serves as an attribute for easy grepping; if not provided it defaults to `HDlog`.

## Features

- Dump headers to traefik's log
- Dump TLS connection, with peer certs into traefik's log
- Customise log prefix for easier grepping or spotting it live

## Deployment

```bash
# static config via CLI argument
--experimental.plugins.headerdump.modulename=github.com/jaybubs/headerdump
--experimental.plugins.headerdump.version=v0.2.0
```

```yaml
# k8s middleware

apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: headerdump
spec:
  plugin:
    headerdump:
      Prefix: "HDlog"
      TLS: true
```

There are no tests, no cicd, just shove it in your route and get investigating.
