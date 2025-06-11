#!/bin/bash
set -e

# Configuration
PROJECT_ID=${GOOGLE_CLOUD_PROJECT:-""}
ZONE="us-central1-c"
INSTANCE_NAME="comfyui-api"
MACHINE_TYPE="n1-standard-4"
GPU_TYPE="nvidia-tesla-t4"
GPU_COUNT=1
IMAGE_FAMILY="ubuntu-2004-lts"
IMAGE_PROJECT="ubuntu-os-cloud"
DISK_SIZE="100GB"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}"
    exit 1
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if gcloud is installed
    if ! command -v gcloud &> /dev/null; then
        error "gcloud CLI is not installed. Please install it first."
    fi
    
    # Check if authenticated
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
        error "Not authenticated with gcloud. Run: gcloud auth login"
    fi
    
    # Check project ID
    if [ -z "$PROJECT_ID" ]; then
        PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
        if [ -z "$PROJECT_ID" ]; then
            error "No GCP project set. Run: gcloud config set project YOUR_PROJECT_ID"
        fi
    fi
    
    log "Using project: $PROJECT_ID"
}

# Enable required APIs
enable_apis() {
    log "Enabling required GCP APIs..."
    gcloud services enable compute.googleapis.com --project=$PROJECT_ID
    gcloud services enable container.googleapis.com --project=$PROJECT_ID
}

# Create VM instance
create_vm() {
    log "Creating VM instance: $INSTANCE_NAME"
    
    # Check if instance already exists
    if gcloud compute instances describe $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID &>/dev/null; then
        warn "Instance $INSTANCE_NAME already exists. Skipping creation."
        return 0
    fi
    
    gcloud compute instances create $INSTANCE_NAME \
        --project=$PROJECT_ID \
        --zone=$ZONE \
        --machine-type=$MACHINE_TYPE \
        --network-interface=network-tier=PREMIUM,subnet=default \
        --maintenance-policy=TERMINATE \
        --provisioning-model=STANDARD \
        --service-account=$(gcloud iam service-accounts list --format="value(email)" --filter="displayName:Compute Engine default service account" --project=$PROJECT_ID) \
        --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
        --accelerator=count=$GPU_COUNT,type=$GPU_TYPE \
        --tags=http-server,https-server,comfyui-api \
        --create-disk=auto-delete=yes,boot=yes,device-name=$INSTANCE_NAME,image=projects/$IMAGE_PROJECT/global/images/family/$IMAGE_FAMILY,mode=rw,size=$DISK_SIZE,type=projects/$PROJECT_ID/zones/$ZONE/diskTypes/pd-standard \
        --no-shielded-secure-boot \
        --shielded-vtpm \
        --shielded-integrity-monitoring \
        --reservation-affinity=any \
        --metadata-from-file startup-script=startup-script.sh
}

# Create firewall rule
create_firewall() {
    log "Creating firewall rule for ComfyUI API..."
    
    # Check if rule already exists
    if gcloud compute firewall-rules describe comfyui-api --project=$PROJECT_ID &>/dev/null; then
        warn "Firewall rule 'comfyui-api' already exists. Skipping creation."
        return 0
    fi
    
    gcloud compute firewall-rules create comfyui-api \
        --project=$PROJECT_ID \
        --direction=INGRESS \
        --priority=1000 \
        --network=default \
        --action=ALLOW \
        --rules=tcp:8188 \
        --source-ranges=0.0.0.0/0 \
        --target-tags=comfyui-api
}

# Wait for VM to be ready
wait_for_vm() {
    log "Waiting for VM to be ready..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID --command="echo 'VM is ready'" &>/dev/null; then
            log "VM is ready and accessible via SSH"
            return 0
        fi
        
        log "Attempt $attempt/$max_attempts - VM not ready yet..."
        sleep 10
        ((attempt++))
    done
    
    error "VM failed to become ready within expected time"
}

# Upload files to VM
upload_files() {
    log "Uploading Docker files to VM..."
    
    # Upload Dockerfile
    gcloud compute scp Dockerfile $INSTANCE_NAME:/tmp/Dockerfile --zone=$ZONE --project=$PROJECT_ID
    
    # Move files to correct location and set permissions
    gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID --command="
        sudo mv /tmp/Dockerfile /opt/comfyui/
        sudo chown root:root /opt/comfyui/Dockerfile
    "
}

# Start ComfyUI service
start_service() {
    log "Starting ComfyUI service..."
    
    gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID --command="
        sudo systemctl start comfyui
        sudo systemctl status comfyui --no-pager
    "
}

# Get VM info
get_vm_info() {
    log "Getting VM information..."
    
    EXTERNAL_IP=$(gcloud compute instances describe $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID --format="get(networkInterfaces[0].accessConfigs[0].natIP)")
    
    echo
    echo "=========================================="
    echo "ComfyUI API Deployment Complete!"
    echo "=========================================="
    echo "VM Name: $INSTANCE_NAME"
    echo "Zone: $ZONE"
    echo "External IP: $EXTERNAL_IP"
    echo "API URL: http://$EXTERNAL_IP:8188"
    echo "API Status: http://$EXTERNAL_IP:8188/system_stats"
    echo
    echo "SSH Access:"
    echo "gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID"
    echo
    echo "View Logs:"
    echo "gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID --command='sudo journalctl -u comfyui -f'"
    echo
    echo "Check Status:"
    echo "gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID --command='sudo comfyui-status'"
    echo "=========================================="
}

# Main execution
main() {
    log "Starting ComfyUI deployment to GCP..."
    
    check_prerequisites
    enable_apis
    create_vm
    create_firewall
    wait_for_vm
    upload_files
    start_service
    get_vm_info
    
    log "Deployment completed successfully!"
}

# Handle script arguments
case "${1:-deploy}" in
    "deploy")
        main
        ;;
    "destroy")
        log "Destroying ComfyUI infrastructure..."
        gcloud compute instances delete $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID --quiet
        gcloud compute firewall-rules delete comfyui-api --project=$PROJECT_ID --quiet
        log "Infrastructure destroyed."
        ;;
    "status")
        EXTERNAL_IP=$(gcloud compute instances describe $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID --format="get(networkInterfaces[0].accessConfigs[0].natIP)" 2>/dev/null || echo "Not found")
        echo "VM Status: $(gcloud compute instances describe $INSTANCE_NAME --zone=$ZONE --project=$PROJECT_ID --format="get(status)" 2>/dev/null || echo "Not found")"
        echo "External IP: $EXTERNAL_IP"
        echo "API URL: http://$EXTERNAL_IP:8188"
        ;;
    *)
        echo "Usage: $0 [deploy|destroy|status]"
        echo "  deploy  - Deploy ComfyUI to GCP (default)"
        echo "  destroy - Destroy all GCP resources"
        echo "  status  - Show current status"
        exit 1
        ;;
esac 