import '@unocss/reset/tailwind.css'
import 'virtual:uno.css'
import './styles/main.css'

import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './App.vue'
import { router } from './router'
import { useAuthStore } from './stores/modules/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)

const authStore = useAuthStore()
authStore.restoreSession()

app.use(router)
app.mount('#app')
