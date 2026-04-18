import { defineConfig, presetIcons, presetTypography, presetUno } from 'unocss'

export default defineConfig({
  presets: [presetUno(), presetTypography(), presetIcons()],
  theme: {
    colors: {
      brand: {
        50: '#eff6ff',
        100: '#dbeafe',
        200: '#bfdbfe',
        300: '#93c5fd',
        400: '#60a5fa',
        500: '#3b82f6',
        600: '#2563eb',
        700: '#1d4ed8',
      },
      ink: {
        50: '#f8fafc',
        100: '#f1f5f9',
        200: '#e2e8f0',
        300: '#cbd5e1',
        500: '#64748b',
        700: '#334155',
        900: '#0f172a',
      },
      steel: {
        100: '#e2e8f0',
        300: '#cbd5e1',
        500: '#64748b',
        700: '#334155',
      },
      mint: {
        100: '#dcfce7',
        500: '#22c55e',
        700: '#15803d',
      },
    },
  },
})
