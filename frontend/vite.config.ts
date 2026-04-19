import path from 'node:path'
import UnoCSS from 'unocss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const proxyTarget = env.VITE_DEV_PROXY_TARGET || 'http://backend:8080'

  return {
    plugins: [vue(), UnoCSS()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      host: '0.0.0.0',
      port: 3000,
      allowedHosts: ['devhub.local', 'localhost', '127.0.0.1'],
      proxy: {
        '/api': {
          target: proxyTarget,
          changeOrigin: true,
          rewrite: (value) => value.replace(/^\/api/, ''),
        },
      },
    },
  }
})
