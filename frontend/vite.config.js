import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    port: 5173,
    open: true,
    cors: true,
    hmr: {
      protocol: 'http',
      host: 'localhost',
      port: 5173
    }
  },
  build: {
    outDir: 'dist',
    rollupOptions: {
      input: {
        main: 'index.html',
        admin: 'admin.html',
        viewer: 'viewer.html'
      }
    }
  }
})
