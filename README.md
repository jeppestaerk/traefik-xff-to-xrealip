# traefik-xff-to-xrealip

A [Traefik](https://traefik.io) plugin that sets the `X-Real-Ip` header based on a **configurable IP address** (by index/depth) in the `X-Forwarded-For` header. By default, it uses the first IP.

## 🔧 What It Does

For incoming requests, this plugin:

- Looks for the `X-Forwarded-For` header.
- Splits the header value by commas to get a list of IP addresses.
- Extracts an IP from this list based on the configured `depth` (index). Defaults to `depth: 0` (the first IP).
- Overwrites `X-Real-Ip` with that value if the depth is valid for the list of IPs.

## 🧪 Example

### Default Behavior (depth: 0)

#### Incoming Request:
```
X-Forwarded-For: 203.0.113.5, 10.0.0.1, 192.168.1.100
```

#### After Plugin (with default or `depth: 0`):
```
X-Real-Ip: 203.0.113.5
```

### Configured Depth (e.g., depth: 1)

#### Incoming Request:
```
X-Forwarded-For: 203.0.113.5, 10.0.0.1, 192.168.1.100
```

#### Middleware Configuration:
```yaml
http:
  middlewares:
    xff2realip-depth1:
      plugin:
        traefik-xff-to-xrealip:
          depth: 1 # Selects the second IP (index 1)
```

#### After Plugin (with `depth: 1`):
```
X-Real-Ip: 10.0.0.1
```

## 🚀 Usage

### 1. Static Traefik Configuration

```yaml
experimental:
  plugins:
    traefik-xff-to-xrealip:
      moduleName: github.com/jeppestaerk/traefik-xff-to-xrealip
      version: v0.0.1
```

### 2. Dynamic Configuration

To use the default depth (0, i.e., the first IP):
```yaml
http:
  middlewares:
    xff2realip-default:
      plugin:
        traefik-xff-to-xrealip: {} # No depth specified, defaults to 0
```

To specify a custom depth (e.g., to select the second IP, index 1):
```yaml
http:
  middlewares:
    xff2realip-custom-depth:
      plugin:
        traefik-xff-to-xrealip:
          depth: 1
```

Apply the middleware to your routers:

```yaml
http:
  routers:
    my-app:
      rule: Host(`myapp.example.com`)
      service: my-app
      middlewares:
        - xff2realip@file
```

## 📦 Installation

This plugin must be built or declared through Traefik's [experimental plugin system](https://doc.traefik.io/traefik/plugins/overview/).  
You can run Traefik with plugins using Docker, Kubernetes, or binary.

## 🧰 Development

```bash
go test ./...
```

GitHub Actions CI is set up to run tests on every commit or pull request.

## 📝 License

MIT © 2025 Jeppe Stærk