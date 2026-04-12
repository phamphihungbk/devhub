<script setup lang="ts">
defineProps<{
  label: string
  tone: 'accent' | 'positive' | 'neutral' | 'warning'
}>()
</script>

<template>
  <span class="status-badge" :data-tone="tone">
    {{ label }}
  </span>
</template>
