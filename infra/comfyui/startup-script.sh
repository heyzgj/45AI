#!/bin/bash
set -e

# Logging function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a /var/log/comfyui-setup.log
}

log "Starting ComfyUI setup on GCP VM..."

# Update system
log "Updating system packages..."
apt-get update && apt-get upgrade -y

# Install Docker
log "Installing Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
systemctl start docker
systemctl enable docker

# Install NVIDIA Container Toolkit
log "Installing NVIDIA Container Toolkit..."
distribution=$(. /etc/os-release;echo $ID$VERSION_ID) \
   && curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | apt-key add - \
   && curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | tee /etc/apt/sources.list.d/nvidia-docker.list

apt-get update
apt-get install -y nvidia-docker2
systemctl restart docker

# Verify NVIDIA Docker setup
log "Verifying NVIDIA Docker setup..."
docker run --rm --gpus all nvidia/cuda:11.8-base-ubuntu20.04 nvidia-smi

# Create ComfyUI directory
mkdir -p /opt/comfyui
cd /opt/comfyui

# Create docker-compose.yml
log "Creating Docker Compose configuration..."
cat > docker-compose.yml << 'EOF'
version: '3.8'
services:
  comfyui:
    build: .
    ports:
      - "8188:8188"
    volumes:
      - ./models:/app/models
      - ./output:/app/output
      - ./temp:/app/temp
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
    restart: unless-stopped
    environment:
      - NVIDIA_VISIBLE_DEVICES=all
      - NVIDIA_DRIVER_CAPABILITIES=compute,utility
EOF

# Copy Dockerfile (will be uploaded separately)
log "Dockerfile should be uploaded to /opt/comfyui/"

# Create systemd service for auto-start
log "Creating systemd service..."
cat > /etc/systemd/system/comfyui.service << 'EOF'
[Unit]
Description=ComfyUI Docker Service
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/comfyui
ExecStart=/usr/local/bin/docker-compose up -d
ExecStop=/usr/local/bin/docker-compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF

# Install docker-compose
log "Installing Docker Compose..."
curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Enable service (will start after Dockerfile is uploaded)
systemctl enable comfyui.service

# Create status endpoint
log "Creating status script..."
cat > /usr/local/bin/comfyui-status << 'EOF'
#!/bin/bash
echo "=== ComfyUI Status ==="
echo "Docker status: $(systemctl is-active docker)"
echo "ComfyUI service: $(systemctl is-active comfyui || echo 'inactive')"
echo "Container status:"
docker ps | grep comfyui || echo "No ComfyUI container running"
echo "GPU status:"
nvidia-smi --query-gpu=name,temperature.gpu,utilization.gpu,memory.used,memory.total --format=csv,noheader,nounits
EOF
chmod +x /usr/local/bin/comfyui-status

log "ComfyUI VM setup completed! Upload Dockerfile and run: systemctl start comfyui"
log "Check status with: comfyui-status"
log "View logs with: journalctl -u comfyui -f" 