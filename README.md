# 45AI - AI Image Generation App

A premium AI-powered image generation application for iOS and WeChat Mini Program, designed for young aesthetically-minded users in China.

## 🎯 Project Status: MVP Complete ✅

**Current Phase**: Production-Ready MVP with Mock Services  
**Last Updated**: January 9, 2025

45AI provides a one-click AI portrait generation experience with professionally designed templates. Users upload a single selfie, select a template, and receive beautiful AI-generated portraits instantly.

### ✅ Implemented Features
- 🎨 **Template-based image generation** - Working with Gemini Flash
- 💳 **Credit-based monetization system** - Complete billing with mock payments
- 🔐 **Authentication system** - WeChat login + development mock
- 📱 **Cross-platform UI** - UniBest + wot-design-uni components
- ✨ **Premium UI with animations** - Modern, responsive design
- 📊 **User management** - Profiles, credits, transaction history
- 🖼️ **Generation history** - View and download past results

## 🏗️ Architecture

### Current (Development)
- **Frontend**: UniBest (uni-app) + wot-design-uni - Vue 3 cross-platform
- **Backend**: Go with Gin framework - RESTful API server
- **Database**: SQLite (dev) / MySQL (production) - User and transaction data
- **AI Service**: Gemini Flash (temporary) → ComfyUI (planned)
- **Authentication**: WeChat OAuth + JWT tokens
- **Payments**: Mock services (WeChat Pay + Apple IAP)

### Production Target
- **Infrastructure**: WeChat Cloud Hosting + GCP
- **AI Service**: Self-hosted ComfyUI on GCP
- **Database**: MySQL 8.0 with proper scaling
- **Storage**: Cloud storage for generated images
- **Payments**: Real WeChat Pay and Apple IAP integration

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

## 🚀 Quick Start (Development)

### Prerequisites
- Go 1.20+
- Node.js 18+
- WeChat Developer Tools (optional)

### Backend Setup
```bash
cd backend
go mod download
# Uses SQLite by default - no additional setup needed
go run cmd/api/main.go
# Server starts on http://localhost:8080
```

### Frontend Setup
```bash
cd frontend
npm install
# For H5 development (recommended)
npm run dev:h5
# For WeChat Mini Program
npm run dev:mp-weixin
```

### Development Login
Use the mock authentication code: `dev_mock_code`  
Mock user: "月来公主" with 100 credits

### Test the System
1. Start backend: `cd backend && go run cmd/api/main.go`
2. Start frontend: `cd frontend && npm run dev:h5`
3. Open http://localhost:3000
4. Login with `dev_mock_code`
5. Browse templates and generate images!

## 📋 Project Management

The project uses a file-based task management system. Check `/project/project_status.md` for current progress and active tasks.

## 📚 Documentation

- [**Development Progress**](docs/DEVELOPMENT_PROGRESS.md) - Complete status summary
- [Product Requirements](docs/PRD.md) - Feature specifications
- [Architecture Guide](docs/ARCHITECTURE_GUIDE.md) - System design
- [Technical Specification](docs/TECH_SPEC.md) - Implementation details
- [Style Guide](docs/STYLE_GUIDE.md) - UI/UX guidelines
- [Test Plan](docs/TEST_PLAN.md) - Testing strategy

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