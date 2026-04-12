<script setup lang="ts">
import type { GlobalThemeOverrides } from 'naive-ui'
import { NConfigProvider, NDialogProvider, NGlobalStyle, NLoadingBarProvider, NMessageProvider, darkTheme } from 'naive-ui'

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#2563eb',
    primaryColorHover: '#1d4ed8',
    primaryColorPressed: '#1e40af',
    successColor: '#0f766e',
    warningColor: '#d97706',
    errorColor: '#dc2626',
    borderRadius: '16px',
    fontFamily: '"Inter", "Avenir Next", system-ui, sans-serif',
  },
  Card: {
    borderRadius: '24px',
  },
  Button: {
    borderRadiusMedium: '14px',
  },
}
</script>

<template>
  <NConfigProvider :theme="darkTheme" :theme-overrides="themeOverrides">
    <NLoadingBarProvider>
      <NDialogProvider>
        <NMessageProvider>
          <NGlobalStyle />
          <NuxtLayout>
            <NuxtPage />
          </NuxtLayout>
        </NMessageProvider>
      </NDialogProvider>
    </NLoadingBarProvider>
  </NConfigProvider>
</template>
