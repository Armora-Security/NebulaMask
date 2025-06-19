# NebulaMask
The Load Balancer in the Road Warrior Muscle

# 🛡️ Armora NebulaMask

> Modern Load Balancer with HTTP-to-HTTPS Redirect, Health Checks & Round Robin Distribution  
> 🔒 Copyright Protected – Commercial Use Requires License  

## 📌 Introduction

**Armora NebulaMask** is a high-performance load balancer written in Go (Golang), designed to efficiently route traffic across multiple backend servers. This tool supports:

- **Automatic HTTP to HTTPS redirection**
- **Round Robin** traffic distribution
- Periodic **health checks** on backend services
- Support for up to **10 backend servers**
- Easy configuration via YAML file
- SSL/TLS termination support

It's ideal for environments where simplicity, reliability, and performance are key.

---

## ⚙️ How It Works

### Architecture Overview
Client --> [Port 80 (HTTP)] --(Redirect to HTTPS)--> [Port 443 (HTTPS)] | ↓ [Backend Servers] (Up to 10 backend nodes with health monitoring)


### Core Features

- **Load Balancing**: Distributes incoming requests using round-robin algorithm.
- **Health Monitoring**: Pings each backend at configurable intervals to determine availability.
- **Secure by Default**: Built-in support for TLS termination using provided certificates.
- **Lightweight**: Minimal dependencies and resource-efficient.

---

## 🧪 Requirements

- Go 1.22+ (for building from source)
- Valid SSL certificate (`cert.pem` and `key.pem`)
- Backend servers running and responding to `/health` endpoint

---

## 📁 Configuration Example (`config.yaml`)

```yaml
server:
  http_port: 80
  https_port: 443

ssl:
  enabled: true
  cert_file: "cert.pem"
  key_file: "key.pem"

backends:
  - name: "server-1"
    url: "http://127.0.0.1:3001"
    health_check_path: "/health"
    interval_sec: 5

  - name: "server-2"
    url: "http://127.0.0.1:3002"
    health_check_path: "/health"
    interval_sec: 5

  - name: "server-3"
    url: "http://127.0.0.1:3003"
    health_check_path: "/health"
    interval_sec: 5

  - name: "server-4"
    url: "http://127.0.0.1:3004"
    health_check_path: "/health"
    interval_sec: 5

  - name: "server-5"
    url: "http://127.0.0.1:3005"
    health_check_path: "/health"
    interval_sec: 5

  - name: "server-6"
    url: "http://127.0.0.1:3006"
    health_check_path: "/health"
    interval_sec: 5

  - name: "server-7"
    url: "http://127.0.0.1:3007"
    health_check_path: "/health"
    interval_sec: 5

  - name: "server-8"
    url: "http://127.0.0.1:3008"
    health_check_path: "/health"
    interval_sec: 5

  - name: "server-9"
    url: "http://127.0.0.1:3009"
    health_check_path: "/health"
    interval_sec: 5

  - name: "server-10"
    url: "http://127.0.0.1:3010"
    health_check_path: "/health"
    interval_sec: 5

```

## 🛠️ Build Instructions

To compile:

```bash

go build -o nebulamask
```
 For cross-compilation:

```bash

GOOS=linux GOARCH=amd64 go build -o nebulamask
```

 Or use Docker:

```bash

docker build -t nebulamask .
docker run -p 80:80 -p 443:443 nebulamask
```
## ⚠️ Important Notices

❗ This software is NOT free for commercial use. You must obtain a valid license from the owner(s) to use it in production or business environments. 

 ✅ Only authorized contributors may use this software freely for personal, non-commercial purposes. Redistribution or resale is strictly prohibited without prior written permission. 

 ⚠️ USE AT YOUR OWN RISK. The developers and maintainers of Armora NebulaMask assume no responsibility for any damage, data loss, downtime, or other consequences resulting from the use of this software. You bear all risks associated with its use. 

## 📄 License

See LICENSE.txt for full license terms.
