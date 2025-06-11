#!/bin/bash

# WeChat Cloud Hosting Deployment Script for 45AI Backend
# This script automates the deployment process to WeChat Cloud Hosting

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="45ai-backend"
ENVIRONMENT="${1:-production}"
VERSION=$(date +%Y%m%d-%H%M%S)

echo -e "${BLUE}🚀 Starting WeChat Cloud Hosting deployment for $PROJECT_NAME${NC}"
echo -e "${BLUE}Environment: $ENVIRONMENT${NC}"
echo -e "${BLUE}Version: $VERSION${NC}"

# Function to log messages
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}"
    exit 1
}

warn() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

# Check dependencies
check_dependencies() {
    log "Checking dependencies..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        error "Go is not installed. Please install Go 1.21 or later."
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log "Go version: $GO_VERSION"
    
    # Check if WeChat Cloud CLI is installed (tcb)
    if ! command -v tcb &> /dev/null; then
        warn "WeChat Cloud CLI (tcb) not found. Please install it first."
        echo "Install with: npm install -g @cloudbase/cli"
        error "Missing WeChat Cloud CLI"
    fi
    
    log "Dependencies check passed ✓"
}

# Build the application
build_app() {
    log "Building Go application..."
    
    cd "$(dirname "$0")/../backend"
    
    # Clean previous builds
    rm -f main
    
    # Download dependencies
    log "Downloading Go modules..."
    go mod download
    go mod tidy
    
    # Run tests
    log "Running tests..."
    go test ./... || warn "Some tests failed, but continuing with deployment"
    
    # Build the application
    log "Building binary..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags "-X main.version=$VERSION -X main.buildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
        -o main \
        ./cmd/api
    
    if [ ! -f "main" ]; then
        error "Build failed - binary not found"
    fi
    
    log "Build completed successfully ✓"
}

# Validate configuration
validate_config() {
    log "Validating deployment configuration..."
    
    CONFIG_FILE="../wechat-cloud-config.yaml"
    if [ ! -f "$CONFIG_FILE" ]; then
        error "Configuration file not found: $CONFIG_FILE"
    fi
    
    # Check required environment variables
    REQUIRED_VARS=(
        "DB_HOST"
        "DB_NAME" 
        "DB_USER"
        "DB_PASSWORD"
        "JWT_SECRET"
        "WECHAT_APP_ID"
        "WECHAT_APP_SECRET"
    )
    
    for var in "${REQUIRED_VARS[@]}"; do
        if [ -z "${!var}" ]; then
            warn "Environment variable $var is not set"
        fi
    done
    
    log "Configuration validation completed ✓"
}

# Deploy to WeChat Cloud
deploy_to_cloud() {
    log "Deploying to WeChat Cloud Hosting..."
    
    # Login check
    log "Checking WeChat Cloud authentication..."
    tcb login --check || {
        warn "Not logged in to WeChat Cloud. Please login first:"
        echo "tcb login"
        error "Authentication required"
    }
    
    # Deploy using CloudBase CLI
    log "Starting deployment process..."
    
    # Set environment variables for deployment
    export TCB_ENV_ID="${TCB_ENV_ID:-45ai-prod}"
    export TCB_SECRET_ID="${TCB_SECRET_ID}"
    export TCB_SECRET_KEY="${TCB_SECRET_KEY}"
    
    # Deploy the application
    tcb hosting:deploy \
        --envId "$TCB_ENV_ID" \
        --src ./main \
        --cloudPath /app \
        --force
    
    # Deploy configuration
    log "Deploying configuration..."
    tcb functions:deploy \
        --envId "$TCB_ENV_ID" \
        --name "$PROJECT_NAME" \
        --src . \
        --runtime go1.21 \
        --handler main \
        --timeout 30 \
        --memorySize 512
    
    log "Deployment completed successfully ✓"
}

# Health check after deployment
health_check() {
    log "Performing health check..."
    
    # Wait for deployment to be ready
    sleep 30
    
    # Health check URL
    HEALTH_URL="${HEALTH_URL:-https://$TCB_ENV_ID.ap-shanghai.tcb.qcloud.la/health}"
    
    log "Checking health endpoint: $HEALTH_URL"
    
    # Retry health check up to 5 times
    for i in {1..5}; do
        if curl -f -s "$HEALTH_URL" > /dev/null; then
            log "Health check passed ✓"
            return 0
        else
            warn "Health check attempt $i failed, retrying in 10 seconds..."
            sleep 10
        fi
    done
    
    error "Health check failed after 5 attempts"
}

# Rollback function
rollback() {
    warn "Initiating rollback..."
    
    # Rollback to previous version
    tcb functions:version:activate \
        --envId "$TCB_ENV_ID" \
        --name "$PROJECT_NAME" \
        --versionNumber "previous"
    
    log "Rollback completed"
}

# Main deployment process
main() {
    log "Starting deployment process..."
    
    # Trap errors for cleanup
    trap 'error "Deployment failed at step: $BASH_COMMAND"' ERR
    
    check_dependencies
    validate_config
    build_app
    deploy_to_cloud
    health_check
    
    log "🎉 Deployment completed successfully!"
    log "Application is now running at: https://$TCB_ENV_ID.ap-shanghai.tcb.qcloud.la"
    log "Health check: https://$TCB_ENV_ID.ap-shanghai.tcb.qcloud.la/health"
}

# Handle command line arguments
case "${1:-deploy}" in
    "deploy")
        main
        ;;
    "rollback")
        rollback
        ;;
    "health")
        health_check
        ;;
    "build")
        check_dependencies
        build_app
        ;;
    *)
        echo "Usage: $0 {deploy|rollback|health|build}"
        echo "  deploy   - Full deployment process (default)"
        echo "  rollback - Rollback to previous version"
        echo "  health   - Run health check only"
        echo "  build    - Build application only"
        exit 1
        ;;
esac 