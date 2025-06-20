import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

import 'uno.css'
import 'wot-design-uni/dist/style.css'
import './styles/index.scss'

export function createApp() {
  const app = createSSRApp(App)
  const pinia = createPinia()
  
  app.use(pinia)
  
  return {
    app,
    pinia,
  }
} 