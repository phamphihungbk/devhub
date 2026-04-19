<script setup lang="ts">
defineProps<{
  eyebrow: string
  title: string
  description: string
}>()
</script>

<template>
  <header class="mb-6 flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
    <div>
      <p class="text-xs font-700 uppercase tracking-0.24em text-brand-700">{{ eyebrow }}</p>
      <h2 class="mt-2 text-3xl font-800 tracking-tight text-ink-900">{{ title }}</h2>
      <p class="mt-2 max-w-2xl text-sm leading-6 text-[var(--app-text-muted)]">{{ description }}</p>
    </div>
    <slot />
  </header>
</template>
