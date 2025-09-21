# OpenGate

<div align="center">
  <img src="resources/images/golang.png" alt="OpenGate Logo" width="200"/>
  
  **A High-Performance API Gateway Built with Go**
  
  [![Go Version](https://img.shields.io/badge/Go-1.23.3-blue.svg)](https://golang.org)
  [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
  [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
</div>

## 🚀 Overview

OpenGate is a modern, lightweight API Gateway built in Go that provides intelligent request routing, authentication, and middleware management for microservices architectures. It offers dynamic route configuration, built-in authentication strategies, and high-performance request forwarding.

## ✨ Key Features

- **🔀 Dynamic Route Management**: File-based route configuration with hot-reload capabilities
- **🔐 Flexible Authentication**: Support for multiple authentication strategies including OpenAuth
- **⚡ High Performance**: Built on Gin framework for optimal request handling
- **🔧 Middleware Support**: Extensible middleware pipeline with CORS, logging, and custom middleware
- **📊 Health Monitoring**: Built-in health checks and ping endpoints
- **🐳 Container Ready**: Docker support for easy deployment
- **📁 Multiple Backends**: Support for local file-based and remote configuration sources
- **⏱️ Configurable Timeouts**: Per-route timeout configuration
- **🔄 Hot Configuration Reload**: Automatic detection and reload of configuration changes

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Client        │    │   OpenGate       │    │   Backend       │
│   Application   │────│   API Gateway    │────│   Services      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Configuration   │
                    │  Management      │
                    └──────────────────┘
```

OpenGate acts as a reverse proxy and load balancer, sitting between client applications and backend services, providing:
- Request routing based on path prefixes
- Authentication and authorization
- Request/response transformation
- Traffic management and monitoring

## 📦 Installation

### Prerequisites

- Go 1.23.3 or higher
- Docker (optional, for containerized deployment)

### From Source

```bash
# Clone the repository
git clone https://github.com/gofreego/opengate.git
cd opengate

# Install dependencies
go mod tidy

# Build the application
make build

# Run the application
./application -env=dev -path=./
```

### Using Docker

```bash
# Build Docker image
make docker

# Run with Docker
make docker-run
```

### Using Go Install

```bash
# Install required tools
make install

# Run directly
make run
```

## ⚙️ Configuration

OpenGate uses YAML configuration files for setup. The main configuration file (`dev.yaml`) defines:

### Server Configuration

```yaml
Server:
  Port: 8080
  GinMode: debug
  ReadTimeout: 30s
  WriteTimeout: 30s
  IdleTimeout: 120s
  MaxHeaderBytes: 1048576
  EnableCors: true
```

### Repository Configuration

```yaml
Repository:
  Name: Local  # or "OpenAuth" for remote configuration
  Local: 
    RoutesFolderPath: ./resources/configs/routes/
  OpenAuth:
    endpoint: localhost:8086
    username: admin
    password: admin123
    timeout: 5s
    tls: false
```

### Service Configuration

```yaml
Service:
  ChangeDetector:
    RouteUpdateInterval: 30  # seconds
  Auth:
    Name: "OpenAuth"
```

## 🛣️ Route Configuration

Routes are defined in YAML files within the routes directory. Each service gets its own configuration file:

```yaml
# resources/configs/routes/example-service.yaml
Name: example-service
PathPrefix: /api/v1/users
TargetURL: http://localhost:3001
StripPrefix: false

Authentication:
  Required: true
  Except:
    - Path: "/api/v1/users/health"
      Methods: ["GET"]
    - Path: "/api/v1/users/metrics"
      Methods: ["GET", "POST"]

Middleware:
  - cors
  - logging

Timeout: 30s
```

### Route Configuration Options

| Field | Type | Description |
|-------|------|-------------|
| `Name` | string | Service identifier for logging and management |
| `PathPrefix` | string | URL path prefix that triggers this route |
| `TargetURL` | string | Backend service URL where requests are forwarded |
| `StripPrefix` | boolean | Whether to remove the path prefix before forwarding |
| `Authentication.Required` | boolean | Whether authentication is required |
| `Authentication.Except` | array | Paths/methods exempt from auth requirements |
| `Middleware` | array | Ordered list of middleware to apply |
| `Timeout` | duration | Request timeout for this route |

## 🔐 Authentication

OpenGate supports multiple authentication strategies:

### OpenAuth Integration

```yaml
Service:
  Auth:
    Name: "OpenAuth"
    
Repository:
  OpenAuth:
    endpoint: auth-service:8086
    username: gateway-user
    password: secure-password
    timeout: 5s
    tls: true
```

### Route-Level Authentication

```yaml
Authentication:
  Required: true
  Except:
    - Path: "/health"
      Methods: ["GET"]
    - Path: "/public/*"
      Methods: ["GET", "POST"]
```

## 🚀 Getting Started

### 1. Basic Setup

```bash
# Clone and build
git clone https://github.com/gofreego/opengate.git
cd opengate
make build
```

### 2. Configure Your First Route

Create a route configuration file:

```yaml
# resources/configs/routes/my-service.yaml
Name: my-service
PathPrefix: /api/my-service
TargetURL: http://localhost:3000
StripPrefix: true
Authentication:
  Required: false
Middleware:
  - cors
Timeout: 30s
```

### 3. Start the Gateway

```bash
./application -env=dev -path=./
```

### 4. Test Your Route

```bash
# Health check
curl http://localhost:8080/ping

# Test your service route
curl http://localhost:8080/api/my-service/endpoint
```

## 🛠️ Development

### Project Structure

```
opengate/
├── cmd/
│   └── http_server/          # HTTP server implementation
├── internal/
│   ├── configs/              # Configuration management
│   ├── constants/            # Application constants
│   ├── models/               # Data models
│   ├── repository/           # Data access layer
│   └── service/              # Business logic
│       ├── auth/             # Authentication services
│       ├── change_detector/  # Configuration change detection
│       └── route_manager/    # Route management
├── pkg/
│   └── utils/                # Utility packages
├── resources/
│   ├── configs/routes/       # Route configurations
│   └── images/               # Static assets
└── test/                     # Test files
```

### Building and Testing

```bash
# Run tests
make test

# Build for Linux
make build-linux

# Clean build artifacts
make clean

# Run development server
make run
```

### Adding New Routes

1. Create a new YAML file in `resources/configs/routes/`
2. Define your route configuration
3. OpenGate will automatically detect and load the new route
4. Test the route with your preferred HTTP client

### Custom Middleware

OpenGate supports extensible middleware. To add custom middleware:

1. Implement the middleware in the service layer
2. Register it in the middleware registry
3. Reference it in your route configuration

## 📊 Monitoring and Health Checks

### Health Endpoints

- `GET /ping` - Basic health check
- `GET /health` - Detailed health status (if configured)

### Logging

OpenGate provides structured logging with configurable levels:

```yaml
Logger:
  AppName: opengate
  Build: dev
  Level: debug  # debug, info, warn, error
```

## 🐳 Docker Deployment

### Using the Provided Dockerfile

```bash
# Build image
docker build -t opengate .

# Run container
docker run -d \
  --name opengate \
  -p 8080:8080 \
  -v $(pwd)/resources/configs:/app/resources/configs \
  opengate
```

### Docker Compose Example

```yaml
version: '3.8'
services:
  opengate:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./resources/configs:/app/resources/configs
      - ./dev.yaml:/app/dev.yaml
    environment:
      - ENV=production
```

## ⚡ Performance

OpenGate is designed for high performance:

- **Gin Framework**: Fast HTTP routing and middleware
- **Connection Pooling**: Efficient backend connections
- **Configurable Timeouts**: Prevent hanging requests
- **Hot Reload**: Zero-downtime configuration updates
- **Lightweight**: Minimal resource footprint

### Performance Tuning

```yaml
Server:
  ReadTimeout: 30s      # Adjust based on your needs
  WriteTimeout: 30s     # Balance between performance and reliability
  IdleTimeout: 120s     # Connection keep-alive duration
  MaxHeaderBytes: 1048576  # Maximum header size
```

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new functionality
- Update documentation for API changes
- Ensure all tests pass before submitting

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: [Wiki](https://github.com/gofreego/opengate/wiki)
- **Issues**: [GitHub Issues](https://github.com/gofreego/opengate/issues)
- **Discussions**: [GitHub Discussions](https://github.com/gofreego/opengate/discussions)

## 🗺️ Roadmap

- [ ] Rate limiting and throttling
- [ ] WebSocket support
- [ ] Circuit breaker pattern
- [ ] Metrics and observability improvements
- [ ] gRPC gateway support
- [ ] Load balancing strategies
- [ ] Admin API for runtime configuration

## 🙏 Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- Powered by [GoUtils](https://github.com/gofreego/goutils) - Utility libraries
- Authentication via [OpenAuth](https://github.com/gofreego/openauth) - Authentication service

---

<div align="center">
  Made with ❤️ by the GoFreego Team
</div>
