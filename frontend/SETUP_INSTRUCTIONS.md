# UniBest Frontend Setup Instructions

## Prerequisites

1. **Install pnpm** (required for UniBest):
   ```bash
   # Option 1: Using npm with sudo
   sudo npm install -g pnpm
   
   # Option 2: Using brew (recommended for macOS)
   brew install pnpm
   
   # Option 3: Using curl
   curl -fsSL https://get.pnpm.io/install.sh | sh -
   ```

2. **Create UniBest project**:
   ```bash
   # Navigate to frontend directory
   cd frontend
   
   # Create UniBest project (interactive)
   pnpm create unibest
   
   # Follow the prompts:
   # - Project name: 45ai-frontend
   # - Select features: TypeScript, UnoCSS, Pinia
   # - UI Library: wot-ui (for WeChat compatibility)
   ```

3. **Move files to correct location**:
   ```bash
   # Move generated files from 45ai-frontend to current directory
   mv 45ai-frontend/* .
   mv 45ai-frontend/.* . 2>/dev/null || true
   rmdir 45ai-frontend
   ```

4. **Install dependencies**:
   ```bash
   pnpm install
   ```

## Development Commands

- **H5 (Web)**: `pnpm dev:h5`
- **WeChat Mini Program**: `pnpm dev:mp-weixin`
- **iOS App**: `pnpm dev:app-plus`

## Build Commands

- **H5 (Web)**: `pnpm build:h5`
- **WeChat Mini Program**: `pnpm build:mp-weixin`
- **iOS App**: `pnpm build:app-plus`

## Notes

- UniBest includes Vue 3, TypeScript, Vite 5, UnoCSS, and wot-ui by default
- The project structure is optimized for cross-platform development
- All our custom styles should override the default wot-ui theme to match our Style Guide 