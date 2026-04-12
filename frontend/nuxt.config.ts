export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },
  modules: ['nuxtjs-naive-ui'],
  css: ['~/assets/css/main.css'],
  app: {
    head: {
      title: 'DevHub',
      meta: [
        {
          name: 'viewport',
          content: 'width=device-width, initial-scale=1',
        },
        {
          name: 'description',
          content:
            'DevHub control plane for projects, plugins, deployments, and platform workflows.',
        },
      ],
    },
  },
  runtimeConfig: {
    backendBaseUrl: process.env.NUXT_BACKEND_BASE_URL || 'http://localhost:8080',
    public: {
      appName: process.env.NUXT_PUBLIC_APP_NAME || 'DevHub',
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || '/api',
    },
  },
})
