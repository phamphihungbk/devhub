import type { GlobalThemeOverrides } from 'naive-ui'
import { defineStore } from 'pinia'

import { SetupStoreId } from '@/enum'
import { themeOverrides } from '@/theme/settings'

function createThemeOverrides(): GlobalThemeOverrides {
  return structuredClone(themeOverrides)
}

export interface ThemeState {
  darkMode: boolean
  naiveTheme: GlobalThemeOverrides
}

export const useThemeStore = defineStore(SetupStoreId.Theme, {
  state: (): ThemeState => ({
    darkMode: false,
    naiveTheme: createThemeOverrides(),
  }),
  getters: {
    isDarkMode(state) {
      return state.darkMode
    },
  },
  actions: {
    setDarkMode(enabled: boolean) {
      this.darkMode = enabled
    },
    toggleDarkMode() {
      this.darkMode = !this.darkMode
    },
    setNaiveTheme(overrides: GlobalThemeOverrides) {
      this.naiveTheme = structuredClone(overrides)
    },
    resetNaiveTheme() {
      this.naiveTheme = createThemeOverrides()
    },
  },
})
