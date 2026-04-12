<script setup lang="ts">
defineProps<{
  title: string
  subtitle: string
  rows: Array<{ id: string }>
  emptyMessage: string
}>()
</script>

<template>
  <section class="resource-table">
    <div class="section-heading">
      <div>
        <p class="section-heading__eyebrow">{{ subtitle }}</p>
        <h2>{{ title }}</h2>
      </div>
      <span class="section-heading__meta">{{ rows.length }}</span>
    </div>
    <div v-if="rows.length === 0" class="empty-state">
      {{ emptyMessage }}
    </div>
    <div v-else class="table-list">
      <article v-for="row in rows" :key="row.id" class="table-row">
        <slot :row="row" />
      </article>
    </div>
  </section>
</template>
