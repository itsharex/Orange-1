import './assets/liquid-glass.css'
import './assets/main.css'
import 'remixicon/fonts/remixicon.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { setAuthLogout } from '@/api'
import { useAuthStore } from '@/stores/auth'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 设置 API 登出回调
setAuthLogout(() => {
  const authStore = useAuthStore()
  authStore.logout()
  router.push('/login')
})

app.mount('#app')
