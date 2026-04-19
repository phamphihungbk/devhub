<script setup lang="ts">
import { computed } from 'vue';
import {
  NConfigProvider,
  NDialogProvider,
  NGlobalStyle,
  NMessageProvider,
  NNotificationProvider,
  darkTheme
} from 'naive-ui'
import { RouterView } from 'vue-router'
import { useThemeStore } from '@/stores/modules/theme'

defineOptions({
  name: 'App'
});

const themeStore = useThemeStore();

const naiveDarkTheme = computed(() => (themeStore.darkMode ? darkTheme : undefined));
</script>

<template>
  <NConfigProvider 
  :theme-overrides="themeStore.naiveTheme"
  :theme="naiveDarkTheme"
  >
    <NMessageProvider>
      <NNotificationProvider>
        <NDialogProvider>
          <NGlobalStyle />
          <RouterView />
        </NDialogProvider>
      </NNotificationProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>
