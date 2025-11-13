# 🐳 Docker Deployment Guide

## Quick Start

The easiest way to run Willowcal is using Docker Compose:

```bash
# Build and start the container
docker-compose up --build

# Or run in detached mode
docker-compose up -d --build
```

Once started, open your browser to:
**http://localhost:8080**

## What Happens When You Run the Container?

1. ✅ Builds the React frontend
2. ✅ Compiles the Go backend
3. ✅ Starts the WebSocket server on port 8080
4. ✅ Serves the web interface
5. ✅ Ready to accept connections!

## Manual Docker Build

If you prefer to build manually:

```bash
# Build the image
docker build -t willowcal:latest .

# Run the container
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/workspace:/app/workspace \
  --name willowcal \
  willowcal:latest
```

## Configuration

### Environment Variables

- `PORT` - Server port (default: 8080)
- `WORKSPACE_DIR` - Directory for cloned repositories (default: /app/workspace)
- `STATIC_DIR` - Directory for web interface files (default: /app/web/dist)

### Volumes

Mount your local workspace directory to persist repositories:

```bash
docker run -d \
  -p 8080:8080 \
  -v /path/to/your/workspace:/app/workspace \
  willowcal:latest
```

## Accessing the Application

- **Web Interface**: http://localhost:8080
- **WebSocket API**: ws://localhost:8080/ws
- **Health Check**: http://localhost:8080/health

## Stopping the Container

```bash
# Using Docker Compose
docker-compose down

# Or manually
docker stop willowcal
docker rm willowcal
```

## Viewing Logs

```bash
# Using Docker Compose
docker-compose logs -f

# Or manually
docker logs -f willowcal
```

## Development Mode

For development with hot-reload:

```bash
# Start backend
go run ./cmd/willowcal server 8080 ./workspace ./web/dist

# In another terminal, start frontend dev server
cd web
npm install
npm run dev
```

Frontend will be available at http://localhost:3000 with hot-reload.

## Troubleshooting

### Port Already in Use

If port 8080 is already in use, change it:

```bash
docker-compose up -d -e PORT=3000
```

Or edit `docker-compose.yml`:

```yaml
ports:
  - "3000:8080"
```

### Permission Issues with Workspace

Ensure the workspace directory has proper permissions:

```bash
chmod -R 755 ./workspace
```

### Container Won't Start

Check the logs:

```bash
docker-compose logs
```

Verify the health status:

```bash
docker-compose ps
```

## Multi-Architecture Support

The image supports both amd64 and arm64:

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t willowcal:latest .
```

## Production Deployment

For production, consider:

1. Using a reverse proxy (nginx/caddy) with SSL
2. Setting up proper authentication
3. Restricting WebSocket origins
4. Implementing rate limiting
5. Using Docker secrets for sensitive data

Example with nginx:

```nginx
server {
    listen 443 ssl;
    server_name willowcal.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```
