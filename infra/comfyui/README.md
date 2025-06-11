# ComfyUI API Infrastructure on GCP

## Overview
This directory contains the infrastructure configuration for deploying ComfyUI API on Google Cloud Platform with GPU support for the 45AI project.

## Prerequisites
- Google Cloud SDK installed and authenticated
- Docker installed locally
- Active GCP project with billing enabled
- Compute Engine API enabled
- Container Registry API enabled

## Architecture
- **VM Instance:** n1-standard-4 with NVIDIA Tesla T4 GPU
- **OS:** Ubuntu 20.04 LTS with GPU drivers
- **ComfyUI:** Running in Docker container
- **Network:** HTTP API exposed on port 8188
- **Storage:** 100GB persistent disk for models

## Quick Setup
1. Set your GCP project: `gcloud config set project YOUR_PROJECT_ID`
2. Run deployment: `./deploy.sh`
3. Verify API: `curl http://[EXTERNAL_IP]:8188/system_stats`

## Files
- `Dockerfile` - ComfyUI container definition
- `startup-script.sh` - VM initialization script
- `deploy.sh` - Deployment automation
- `terraform/` - Infrastructure as Code
- `models/` - Model download scripts

## Estimated Costs
- VM: ~$200/month (24/7 running)
- Storage: ~$10/month
- Network: ~$5/month
- Total: ~$215/month

## Security Notes
- API is currently unprotected - suitable for development only
- Production deployment should include authentication
- Consider using Cloud Run for auto-scaling in production 