# 45AI Backend API

Go-based API server for the 45AI image generation application.

## Prerequisites

- Go 1.20 or higher
- MySQL 8.0
- Docker (optional, for containerized development)

## Installation

1. **Install Go** (if not already installed):
   ```bash
   # macOS
   brew install go
   
   # Or download from https://go.dev/dl/
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /api/v1/ping` - Test endpoint

## Project Structure

```
backend/
├── cmd/
│   └── api/           # Application entry points
├── internal/          # Private application code
│   ├── handler/       # HTTP handlers (controllers)
│   ├── service/       # Business logic
│   ├── repository/    # Data access layer
│   ├── model/         # Data models
│   └── middleware/    # HTTP middleware
├── pkg/              # Public packages
├── migrations/       # Database migrations
└── configs/          # Configuration files
```

## Development

### Running with hot reload:
```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Run with air
air
```

### Running tests:
```bash
go test ./...
```

### Building for production:
```bash
go build -o bin/api cmd/api/main.go
```

## Deployment

The application is designed to run on WeChat Cloud Hosting. See the deployment guide in `/infra` for details. 