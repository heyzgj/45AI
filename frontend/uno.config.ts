import { defineConfig } from 'unocss'
import presetWeapp from 'unocss-preset-weapp'
import { transformerClass } from 'unocss-preset-weapp/transformer'

export default defineConfig({
  presets: [
    presetWeapp(),
  ],
  transformers: [
    transformerClass(),
  ],
  theme: {
    colors: {
      primary: '#E89B93', // Dusty Rose
      secondary: '#F3D9D7', // Powder Pink
      bg: '#FCFBF9', // Alabaster
      surface: '#FFFFFF', // White
      text: '#4A4A4A', // Charcoal
      'text-subtle': '#9B9B9B', // Stone Grey
    },
    fontFamily: {
      sans: ['PingFang SC', 'Inter', 'sans-serif'],
    },
    animation: {
      keyframes: {
        'fade-in': '{from{opacity:0}to{opacity:1}}',
        'slide-up': '{from{transform:translateY(16px);opacity:0}to{transform:translateY(0);opacity:1}}',
        'pulse-soft': '{0%,100%{opacity:1;transform:scale(1)}50%{opacity:0.8;transform:scale(1.05)}}',
      },
      durations: {
        'fade-in': '400ms',
        'slide-up': '400ms',
        'pulse-soft': '2s',
      },
      timingFns: {
        'custom': 'cubic-bezier(0.6, 0.05, 0.4, 1)',
      },
    },
  },
  shortcuts: {
    // Buttons
    'btn-primary': 'px-6 py-3 bg-primary text-white rounded-full shadow-primary transition-all duration-300 active:scale-97',
    'btn-secondary': 'px-6 py-3 border border-primary text-primary rounded-full transition-all duration-300',
    
    // Cards
    'card': 'bg-surface border border-gray-100 rounded-16px shadow-soft p-4',
    
    // Animations
    'animate-fade-in': 'animate-fade-in animate-duration-400 animate-ease-custom',
    'animate-slide-up': 'animate-slide-up animate-duration-400 animate-ease-custom',
    'animate-pulse-soft': 'animate-pulse-soft animate-duration-2000 animate-ease-custom animate-iteration-infinite',
  },
  rules: [
    // Custom shadow for primary color
    ['shadow-primary', { 'box-shadow': '0 4px 12px rgba(232, 155, 147, 0.3)' }],
    ['shadow-soft', { 'box-shadow': '0 4px 24px rgba(74, 74, 74, 0.08)' }],
    
    // Custom scale for micro-interactions
    ['scale-97', { transform: 'scale(0.97)' }],
    ['scale-102', { transform: 'scale(1.02)' }],
  ],
}) 