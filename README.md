# traefik-xff-to-xrealip

A [Traefik](https://traefik.io) plugin that sets the `X-Real-Ip` header based on the **first IP address** in the `X-Forwarded-For` header.

## 🔧 What It Does

For incoming requests, this plugin:

- Looks for the `X-Forwarded-For` header
- Extracts the first IP in the comma-separated list
- Overwrites `X-Real-Ip` with that value

## 🧪 Example

### Incoming Request:
```
X-Forwarded-For: 203.0.113.5, 10.0.0.1
```

### After Plugin:
```
X-Real-Ip: 203.0.113.5
```

## 🚀 Usage

### 1. Static Traefik Configuration

```yaml
experimental:
  plugins:
    traefik-xff-to-xrealip:
      moduleName: github.com/YOUR_GITHUB_USERNAME/traefik-xff-to-xrealip
      version: v0.1.0
```

### 2. Dynamic Configuration

```yaml
http:
  middlewares:
    xff2realip:
      plugin:
        traefik-xff-to-xrealip: {}
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