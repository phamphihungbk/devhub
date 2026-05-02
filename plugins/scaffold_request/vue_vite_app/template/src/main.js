import { createApp } from 'vue'
import './style.css'

createApp({
  template: `
    <main>
      <h1>[[SERVICE_NAME]]</h1>
      <p>Vue frontend running in [[ENVIRONMENT]].</p>
    </main>
  `,
}).mount('#app')
