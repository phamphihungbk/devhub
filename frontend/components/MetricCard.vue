<script setup lang="ts">
defineProps<{
  label: string
  value: string
  description: string
  tone: 'accent' | 'positive' | 'neutral'
}>()
</script>

<template>
  <article class="metric-card" :data-tone="tone">
    <p class="metric-card__label">{{ label }}</p>
    <h3>{{ value }}</h3>
    <p class="metric-card__description">{{ description }}</p>
  </article>
</template>
