# 45AI Frontend - UniBest Application

A cross-platform application built with UniBest (uni-app) for iOS and WeChat Mini Program.

## 🎨 Design Philosophy

The UI follows a **minimalist, soft feminine aesthetic** with:
- **Colors**: Dusty rose (#E89B93) and powder pink (#F3D9D7) accents
- **Background**: Soft alabaster (#FCFBF9) instead of harsh white
- **Typography**: Clean and generous with ample line-height
- **Animations**: Fluid transitions with custom cubic-bezier easing
- **Interactions**: Subtle micro-interactions for premium feel

## 🚀 Setup

### Prerequisites
- Node.js 18+
- pnpm 8+ (see SETUP_INSTRUCTIONS.md for installation)
- WeChat Developer Tools (for Mini Program)
- Xcode (for iOS development)

### Quick Start
```bash
# Install pnpm if not already installed
npm install -g pnpm

# Install dependencies
pnpm install

# Run for different platforms
pnpm dev:h5           # Web development
pnpm dev:mp-weixin    # WeChat Mini Program
pnpm dev:app-plus     # iOS App
```

## 📁 Project Structure

```
src/
├── pages/              # Application pages
│   ├── index/         # Home page
│   ├── gallery/       # Template gallery
│   ├── generate/      # Generation flow
│   ├── profile/       # User profile
│   └── login/         # Login page
├── components/         # Reusable components
│   ├── TemplateCard/  # Template display card
│   ├── CreditDisplay/ # Credit balance display
│   └── LoadingAnimation/ # Custom loading states
├── stores/            # Pinia state management
│   ├── user.ts       # User state & auth
│   └── templates.ts  # Templates data
├── api/              # API client layer
├── styles/           # Global styles
│   ├── variables.scss # Design tokens
│   └── animations.scss # Animation utilities
└── utils/            # Helper functions
```

## 🎯 Key Features Implemented

1. **Refined Animations**
   - Page transitions with fade + slide
   - Staggered animations for lists
   - Micro-interactions on all interactive elements

2. **Design System**
   - UnoCSS configured with our brand colors
   - Custom animation utilities
   - Consistent spacing using 4px grid

3. **Cross-Platform Support**
   - Optimized for both iOS and WeChat
   - Platform-specific adjustments where needed
   - Consistent experience across platforms

## 🔧 Development Guidelines

### Animation Standards
- Use the custom easing: `cubic-bezier(0.6, 0.05, 0.4, 1)`
- Page transitions: 400ms duration
- Micro-interactions: 300ms duration
- Always maintain >55 FPS

### Component Development
- Follow Vue 3 Composition API
- Use TypeScript for type safety
- Implement loading and error states
- Add proper accessibility attributes

### Style Guidelines
```scss
// Use design tokens
color: var(--color-primary); // #E89B93

// Consistent spacing
padding: 16px; // 4px grid unit

// Soft shadows
box-shadow: 0 4px 24px rgba(74, 74, 74, 0.08);

// Button states
&:active {
  transform: scale(0.97);
}
```

## 🌟 UniBest Features Used

- **UnoCSS**: For utility-first styling
- **Pinia**: State management
- **Auto Import**: Components and composables
- **TypeScript**: Full type support
- **wot-ui**: Base component library (customized)

## 📱 Platform-Specific Notes

### WeChat Mini Program
- Uses native WeChat login SDK
- Follows Mini Program design guidelines
- Optimized for WeChat's rendering engine

### iOS
- Native-like transitions
- Supports iOS-specific gestures
- Optimized for various iPhone sizes

## 🧪 Testing

```bash
# Type checking
pnpm type-check

# Run in development with hot reload
pnpm dev:mp-weixin
```

## 🚀 Deployment

See main project README for deployment instructions.