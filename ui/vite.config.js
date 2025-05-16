import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'node:path';


// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {host: true},
  resolve: {
    alias: {
      $styles: path.resolve('./src/styles'),
      $assets: path.resolve('./src/assets'),
      $lib: path.resolve('./src/lib'),
    }
  }
})
