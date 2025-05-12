<p align="center"><img src="https://github.com/jeppestaerk/traefik-xff-to-xrealip/blob/main/.assets/icon.svg?raw=true" alt="logo" height="96" width="96"></p>

# 🎯 Traefik X-Real-IP from X-Forwarded-For Plugin

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/jeppestaerk/traefik-xff-to-xrealip?sort=semver&color=blue&logo=github)](https://github.com/jeppestaerk/traefik-xff-to-xrealip/releases/latest) 
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/jeppestaerk/traefik-xff-to-xrealip/build_and_test.yml?branch=main&logo=go)](https://github.com/jeppestaerk/traefik-xff-to-xrealip/actions/workflows/build_and_test.yml) 
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/jeppestaerk/traefik-xff-to-xrealip/build_and_test.yml?branch=main&label=test&logo=go)](https://github.com/jeppestaerk/traefik-xff-to-xrealip/actions/workflows/build_and_test.yml) 
[![Traefik Plugin Catalog](https://img.shields.io/badge/traefik_plugin_catalog-traefik--xff--to--xrealip-blue?logo=traefikproxy)](https://plugins.traefik.io/plugins/68205916e4f1c0f6442c2669/x-real-ip-from-x-forwarded-for)

A [Traefik](https://traefik.io) plugin that intelligently sets the `X-Real-Ip` header by selecting a **configurable IP address** (by index/depth) from the `X-Forwarded-For` header. By default, it uses the first IP, ensuring your backend services see the correct client IP, even behind multiple proxies! 🚀

## 🔧 What It Does

For incoming requests, this plugin:

-   🕵️‍♂️ Looks for the `X-Forwarded-For` header.
-   🔪 Splits the header value by commas to get a list of IP addresses.
-   🎯 Extracts an IP from this list based on the configured `depth` (index). Defaults to `depth: 0` (the first IP).
-   ✍️ Overwrites `X-Real-Ip` with that value if the depth is valid for the list of IPs.

## 🚀 Usage

### Static configuration

> Add to [Static configuration](https://doc.traefik.io/traefik/reference/static-configuration/overview/)

#### Plug In Configuration

Ensure you are using the [latest version](https://github.com/jeppestaerk/traefik-xff-to-xrealip/releases)

```yaml
## Static configuration
experimental:
  plugins:
    traefik-xff-to-xrealip:
      moduleName: github.com/jeppestaerk/traefik-xff-to-xrealip
      version: v0.1.3
```

#### Entry Points Configuration:

Remember to add your proxy IPs to the `forwardedHeaders.trustedIPs` entryPoint configuration in Traefik. Without this, Traefik won't trust the X-Forwarded-For header from your proxies, and this plugin won't work properly.

```yaml
## Static configuration
entryPoints:
  web:
    address: ":80"
    forwardedHeaders:
      trustedIPs:
        - "127.0.0.1/32"
        - "192.168.1.5" # Your proxy IP eg. the ip of the machine running Cloudflare tunnel
        - "172.16.0.0/16" # trust everytihing from docker eg. if running Cloudflare tunnel in docker container
  websecure:
    address: ":443"
    forwardedHeaders:
      trustedIPs:
        - "127.0.0.1/32"
        - "192.168.1.5" # Your proxy IP eg. the ip of the machine running Cloudflare tunnel
        - "172.16.0.0/16" # trust everytihing from docker eg. if running Cloudflare tunnel in docker container
```

Ensure you add `forwardedHeaders.trustedIPs` to all your entryPoints, especially if you redirect HTTP to HTTPS.

### Dynamic configuration

> Add to [dynamic configuration](https://doc.traefik.io/traefik/reference/dynamic-configuration/file/)

#### Middleware Configuration:

```yaml
## Dynamic configuration
http:
  middlewares:
    xff2realip: # Name your middleware instance
      plugin:
        traefik-xff-to-xrealip: {} # Default depth (0)
        # or for a custom depth:
        # traefik-xff-to-xrealip:
        #   depth: 1
```

#### Router Configuration:

```yaml
## Dynamic configuration
http:
  routers:
    my-app:
      rule: Host(`myapp.example.com`)
      service: my-app
      middlewares:
        - xff2realip@file
```

## 🧪 Examples

### Default Behavior (depth: 0)

#### Incoming Request:
```
X-Forwarded-For: 203.0.113.5, 10.0.0.1, 192.168.1.100
```

#### Middleware Configuration:
```yaml
http:
  middlewares:
    xff2realip:
      plugin:
        traefik-xff-to-xrealip: {}
```

#### After Plugin (with default or `depth: 0`):
```
X-Real-Ip: 203.0.113.5
```

### Configured Depth (e.g., `depth: 1`)

#### Incoming Request:
```
X-Forwarded-For: 203.0.113.5, 10.0.0.1, 192.168.1.100
```

#### Middleware Configuration:
```yaml
http:
  middlewares:
    xff2realip:
      plugin:
        traefik-xff-to-xrealip:
          depth: 1 # Selects the second IP (index 1)
```

#### After Plugin (with `depth: 1`):
```
X-Real-Ip: 10.0.0.1
```

## 📦 Installation

This plugin must be built or declared through Traefik's [experimental plugin system](https://doc.traefik.io/traefik/plugins/overview/).  
You can run Traefik with plugins using Docker, Kubernetes, or binary.

## 🧰 Development

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jeppestaerk/traefik-xff-to-xrealip?color=green)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeppestaerk/traefik-xff-to-xrealip)](https://goreportcard.com/report/github.com/jeppestaerk/traefik-xff-to-xrealip)

```bash
go test ./...
```

GitHub Actions CI is set up to run build and tests on every commit or pull request.

## 📝 License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

MIT © 2025 Jeppe Stærk