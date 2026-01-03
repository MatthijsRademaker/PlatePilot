import { defineConfig, presetWind } from 'unocss'

export default defineConfig({
  // Use Wind preset (Tailwind-compatible) with tw- prefix to avoid Quasar conflicts
  presets: [presetWind({ prefix: 'tw-' })],
  rules: [],
})
