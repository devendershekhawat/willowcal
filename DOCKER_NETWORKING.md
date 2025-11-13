# 🌐 Docker Networking for Services

## The Problem

When services run inside the Docker container, they're isolated from your host machine. If a service starts on port 3000 inside the container, you can't access it at `localhost:3000` from your browser.

## Solutions

### ✅ **Option 1: Host Network Mode (Recommended for Development)**

**File: `docker-compose.yml` (default)**

Uses host networking - services bind directly to your host's network interfaces.

```yaml
services:
  willowcal:
    network_mode: "host"
```

**Pros:**
- ✅ All service ports automatically accessible on localhost
- ✅ No need to pre-configure ports
- ✅ Services work exactly as if running natively
- ✅ Dynamic port allocation works

**Cons:**
- ❌ Only works on Linux (Docker Desktop on Mac/Windows doesn't support it properly)
- ❌ Less network isolation
- ❌ Port conflicts with host services

**Usage:**
```bash
docker-compose up
```

Your services are now accessible:
- Willowcal: http://localhost:8080
- Express API: http://localhost:3000 (if that's what your service uses)
- Any service on any port: accessible directly

---

### ✅ **Option 2: Manual Port Mapping**

**File: `docker-compose.ports.yml`**

Explicitly map container ports to host ports.

```yaml
services:
  willowcal:
    ports:
      - "8080:8080"
      - "3000:3000"
      - "5173:5173"
```

**Pros:**
- ✅ Works on all platforms (Linux, Mac, Windows)
- ✅ Better security/isolation
- ✅ Explicit control over exposed ports

**Cons:**
- ❌ Must know ports in advance
- ❌ Need to update docker-compose.yml for new services
- ❌ Port conflicts need manual resolution

**Usage:**
```bash
# Edit docker-compose.ports.yml to add your service ports
docker-compose -f docker-compose.ports.yml up
```

---

### ✅ **Option 3: Run Without Docker (Best for Development)**

Run willowcal directly on your host machine - no container needed!

**Pros:**
- ✅ All services fully accessible
- ✅ No Docker complexity
- ✅ Faster startup
- ✅ Easier debugging

**Setup:**
```bash
# Build frontend
cd web && npm install && npm run build && cd ..

# Build backend
go build -o willowcal ./cmd/willowcal

# Run
./willowcal server 8080 ./workspace ./web/dist
```

**Access:**
- Willowcal: http://localhost:8080
- All services: accessible on their respective ports

---

## Recommendations

### For Development (Linux):
Use **Option 1** (host networking) - it's the easiest and most flexible.

```bash
docker-compose up
```

### For Development (Mac/Windows):
Use **Option 3** (run natively) - host networking doesn't work well on Docker Desktop.

```bash
./willowcal server 8080 ./workspace ./web/dist
```

### For Production/CI:
Use **Option 2** (explicit ports) or deploy services separately with proper networking.

---

## Checking If Services Are Accessible

After starting a service in willowcal, test it:

```bash
# Check if service is running inside container
docker exec willowcal ps aux | grep node

# Test service from host (should work with host networking)
curl http://localhost:3000

# Test service from inside container
docker exec willowcal curl http://localhost:3000
```

---

## Common Issues

### Issue: Can't access service on Mac/Windows with host networking

**Solution:** Host networking doesn't work on Docker Desktop. Use Option 3 (run natively) instead.

### Issue: Port already in use

**Solution:**
- Check what's using the port: `lsof -i :3000`
- Kill the process or use a different port in your service config

### Issue: Service starts but can't connect

**Check:**
1. Service is actually running: Check logs in willowcal terminal
2. Service is listening on correct interface: Should bind to `0.0.0.0`, not `127.0.0.1`
3. Firewall isn't blocking the port

---

## Example Configs

### Express API (port 3000)
```yaml
services:
  - name: backend
    repo: backend-api
    run_command: npm start  # Listens on PORT env var or 3000
```

### Vite Dev Server (port 5173)
```yaml
services:
  - name: frontend
    repo: frontend-app
    run_command: npm run dev -- --host 0.0.0.0  # Important: bind to 0.0.0.0
```

### Python Flask (port 5000)
```yaml
services:
  - name: api
    repo: python-api
    run_command: flask run --host=0.0.0.0 --port=5000
```

**Key:** Always bind to `0.0.0.0` (all interfaces), not `localhost`/`127.0.0.1` (loopback only).

---

## Quick Start

**Linux:**
```bash
docker-compose up
# Services accessible on localhost:PORT
```

**Mac/Windows:**
```bash
# Build
cd web && npm install && npm run build && cd ..
go build -o willowcal ./cmd/willowcal

# Run
./willowcal server 8080 ./workspace ./web/dist
# Services accessible on localhost:PORT
```
