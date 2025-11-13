# 🌿 Willowcal

**Modern service orchestration and management platform with a beautiful web interface**

Willowcal is a powerful tool for orchestrating multiple git repositories and services in parallel, featuring a stunning Linear-inspired web UI for easy management.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)

## ✨ Features

### 🎨 Beautiful Web Interface
- **Linear-inspired design** with dark theme and smooth animations
- **Real-time updates** via WebSocket connection
- **Interactive terminal** with streaming logs
- **Service dashboard** with intuitive controls
- **Config editor** with YAML validation

### 🚀 Core Capabilities
- **Parallel repository cloning** - Clone multiple repos simultaneously
- **Automated setup** - Run setup commands (npm install, etc.)
- **Service management** - Start/stop services with one click
- **Real-time logs** - Stream service logs in real-time
- **Config validation** - YAML validation before execution
- **Smart retry logic** - Automatic retries with intelligent handling

### 🔧 Technical Features
- WebSocket API for real-time communication
- Process management with graceful shutdown
- Individual service control and monitoring
- Config diff detection for updates
- Health check endpoints
- Docker support for easy deployment

## 🚀 Quick Start

### Using Docker (Recommended)

The easiest way to get started:

```bash
# Clone the repository
git clone https://github.com/yourusername/willowcal.git
cd willowcal

# Start with Docker Compose
docker-compose up --build
```

Open your browser to **http://localhost:8080** 🎉

### Manual Installation

#### Prerequisites
- Go 1.24 or later
- Node.js 18 or later
- Git

#### Build from source

```bash
# Install backend dependencies
go mod download

# Build backend
go build -o willowcal ./cmd/willowcal

# Install frontend dependencies and build
cd web
npm install
npm run build
cd ..

# Run the server
./willowcal server 8080 ./workspace ./web/dist
```

Access the web interface at **http://localhost:8080**

## 📖 Usage

### Web Interface

1. **Upload Configuration**
   - Upload a YAML config file or edit inline
   - Validate configuration

2. **Initialize Repositories**
   - Click "Start Initialization"
   - Watch real-time progress in terminal

3. **Manage Services**
   - Start/stop services with play/stop buttons
   - View live logs for each service
   - Monitor service status and PIDs

### CLI Commands

```bash
# Initialize repositories from config
willowcal init config.yaml

# Run services (CLI mode)
willowcal run config.yaml

# Start WebSocket server with web UI
willowcal server [port] [workspace] [static-dir]

# Examples
willowcal server                              # Port 8080, default paths
willowcal server 3000                         # Custom port
willowcal server 3000 ./workspace ./web/dist  # Custom paths
```

### Configuration File

Create a `config.yaml` file:

```yaml
version: "1.0"
workspace_dir: "./workspace"

repositories:
  - name: backend-api
    url: https://github.com/username/backend.git
    path: ./services/backend
    setup_commands:
      - npm install
      - npm run build

  - name: frontend-app
    url: https://github.com/username/frontend.git
    path: ./services/frontend
    setup_commands:
      - npm install

services:
  - name: backend
    repo: backend-api
    run_command: npm start

  - name: frontend
    repo: frontend-app
    run_command: npm run dev
```

## 🎨 Web Interface Features

### Config Management
- Drag-and-drop YAML file upload
- Inline editor with syntax highlighting
- Real-time validation
- Error reporting with detailed messages

### Service Dashboard
- Beautiful service cards with status indicators
- One-click start/stop controls
- PID and uptime monitoring
- Color-coded status badges

### Real-time Terminal
- Streaming logs from all services
- Filter by service name
- Timestamps for each log line
- Color-coded output (stdout/stderr)
- Clear and close controls

### Smooth Animations
- Framer Motion powered transitions
- Glass-morphism effects
- Hover states and interactions
- Page transitions

## 🐳 Docker Deployment

See [DOCKER.md](DOCKER.md) for comprehensive Docker deployment guide.

Quick commands:

```bash
# Start with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Rebuild
docker-compose up --build
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Web Interface                       │
│           (React + Tailwind + Framer)               │
└──────────────────┬──────────────────────────────────┘
                   │ WebSocket
                   ▼
┌─────────────────────────────────────────────────────┐
│              Go WebSocket Server                     │
│   ┌──────────────────────────────────────────┐     │
│   │  Config Handler  │  Service Manager      │     │
│   │  Init Handler    │  Log Broadcaster      │     │
│   └──────────────────────────────────────────┘     │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│              Service Orchestrator                    │
│   ┌──────────────────────────────────────────┐     │
│   │  Git Service  │  Executor Service        │     │
│   │  Worker Pool  │  Process Monitor         │     │
│   └──────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────┘
```

## 🔌 WebSocket API

### Message Types

**Client → Server:**
- `config.upload` - Upload and validate config
- `init.start` - Start initialization
- `service.list` - Get services
- `service.start` - Start a service
- `service.stop` - Stop a service
- `service.status` - Get service status

**Server → Client:**
- `init.progress` - Real-time init logs
- `init.complete` - Init finished
- `service.log` - Service log line
- `service.started` - Service started
- `service.stopped` - Service stopped
- `error` / `success` - Response messages

Example WebSocket message:

```javascript
// Start a service
{
  "type": "service.start",
  "id": "req-123",
  "payload": {
    "service_name": "backend"
  }
}

// Receive log
{
  "type": "service.log",
  "payload": {
    "service_name": "backend",
    "timestamp": "14:23:45",
    "line": "Server started on port 3000",
    "stream": "stdout"
  }
}
```

## 🛠️ Development

### Backend Development

```bash
# Run with hot-reload (using air or similar)
go run ./cmd/willowcal server 8080 ./workspace ./web/dist

# Run tests
go test ./...

# Build
go build -o willowcal ./cmd/willowcal
```

### Frontend Development

```bash
cd web

# Install dependencies
npm install

# Start dev server (http://localhost:3000)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

The dev server proxies WebSocket and API requests to the Go backend at `localhost:8080`.

## 📦 Project Structure

```
willowcal/
├── cmd/willowcal/              # Main entry point
├── internal/
│   ├── api/                    # WebSocket server & handlers
│   ├── commands/               # CLI commands
│   ├── config/                 # Config parsing & validation
│   ├── executor/               # Command execution
│   ├── git/                    # Git operations
│   ├── models/                 # Data models
│   ├── orchestrator/           # Parallel orchestration
│   ├── reporter/               # Progress reporting
│   ├── runner/                 # Service runner
│   └── service/                # Service management
├── web/                        # React frontend
│   ├── src/
│   │   ├── components/         # React components
│   │   │   ├── ConfigEditor.jsx
│   │   │   ├── ServiceCard.jsx
│   │   │   ├── Terminal.jsx
│   │   │   └── InitSection.jsx
│   │   ├── hooks/              # Custom hooks
│   │   │   └── useWebSocket.js
│   │   ├── App.jsx             # Main app component
│   │   ├── main.jsx            # Entry point
│   │   └── index.css           # Styles
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
├── Dockerfile                  # Multi-stage Docker build
├── docker-compose.yml          # Docker orchestration
├── docker-entrypoint.sh        # Container startup
├── go.mod                      # Go dependencies
└── README.md                   # This file
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Inspired by [Linear](https://linear.app) for the beautiful design
- Built with [React](https://react.dev/), [Tailwind CSS](https://tailwindcss.com/), and [Framer Motion](https://www.framer.com/motion/)
- Powered by [Go](https://golang.org/) and [gorilla/websocket](https://github.com/gorilla/websocket)

## 📞 Support

For issues, questions, or contributions, please open an issue on GitHub.

---

**Made with ❤️ using Go and React**
