# 45AI - AI Image Generation App

A premium AI-powered image generation application for iOS and WeChat Mini Program, designed for young aesthetically-minded users in China.

## 🎯 Project Overview

45AI provides a one-click AI portrait generation experience with professionally designed templates. Users upload a single selfie, select a template, and receive beautiful AI-generated portraits instantly.

### Key Features
- 🎨 Template-based image generation
- 💳 Credit-based monetization system ("胶卷")
- 🛡️ Content safety moderation
- 📱 Cross-platform (iOS + WeChat Mini Program)
- ✨ Premium UI with fluid animations

## 🏗️ Architecture

- **Frontend**: UniBest (uni-app) - Vue 3 based cross-platform framework
- **Backend**: Go with Gin framework - High-performance API server
- **Database**: MySQL 8.0 - User and transaction data
- **AI Service**: Self-hosted ComfyUI on GCP - Image generation
- **Infrastructure**: WeChat Cloud Hosting + GCP

## 📁 Project Structure

```
45AI/
├── docs/              # Project documentation
├── project/           # Project management (epics, tasks)
├── frontend/          # UniBest application
├── backend/           # Go API server
├── infra/            # Infrastructure configurations
└── comfyui/          # AI service configurations
```

## 🚀 Getting Started

### Prerequisites
- Go 1.20+
- Node.js 18+
- MySQL 8.0
- WeChat Developer Tools
- Xcode (for iOS development)

### Backend Setup
```bash
cd backend
go mod download
cp .env.example .env
# Configure your .env file
go run cmd/api/main.go
```

### Frontend Setup
```bash
cd frontend
npm install
# For WeChat Mini Program
npm run dev:mp-weixin
# For iOS
npm run dev:app-plus
```

## 📋 Project Management

The project uses a file-based task management system. Check `/project/project_status.md` for current progress and active tasks.

## 📚 Documentation

- [Product Requirements](docs/PRD.md)
- [Architecture Guide](docs/ARCHITECTURE_GUIDE.md)
- [Technical Specification](docs/TECH_SPEC.md)
- [Style Guide](docs/STYLE_GUIDE.md)
- [Test Plan](docs/TEST_PLAN.md)

## 🎨 Design Philosophy

The UI follows a minimalist, soft feminine aesthetic with:
- Dusty rose and powder pink color palette
- Generous whitespace
- Fluid animations with custom easing
- Premium, serene digital experience

## 🔧 Development

This project follows:
- Test-Driven Development (TDD)
- Repository pattern for data access
- Atomic task breakdown
- Continuous integration practices

## 📄 License

Private and Confidential 