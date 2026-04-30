import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '/gateway/',
  resolve: {
    dedupe: ['react', 'react-dom', '@emotion/react', '@emotion/styled', '@mui/material'],
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          // Core React libraries
          if (id.includes('node_modules/react') || id.includes('node_modules/react-dom') || id.includes('node_modules/react-router-dom')) {
            return 'react-vendor';
          }
          // MUI core
          if (id.includes('node_modules/@mui/material') || id.includes('node_modules/@emotion')) {
            return 'mui-core';
          }
          // MUI icons
          if (id.includes('node_modules/@mui/icons-material')) {
            return 'mui-icons';
          }
          // App utilities
          if (id.includes('node_modules/@gofreego/tsutils')) {
            return 'app-utils';
          }
        },
      },
    },
  },
})
