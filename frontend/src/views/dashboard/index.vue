<script setup lang="ts">
import {
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NGrid,
  NGridItem,
  NStatistic,
  NTag,
} from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useDashboardService } from '@/services/dashboard'

const {
  getRoleTagColor,
  loadDashboard,
  loading,
  pluginColumns,
  plugins,
  projectColumns,
  projects,
  stats,
  teamMembers,
} = useDashboardService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Overview"
      title="Platform dashboard"
      description="A higher-signal entry point for the control plane: service inventory, automation registry, and operator visibility in one calmer admin workspace."
    >
      <NButton type="primary" @click="loadDashboard">
        Refresh data
      </NButton>
    </PageHeader>

    <NGrid cols="1 s:3" responsive="screen" :x-gap="16" :y-gap="16">
      <NGridItem v-for="item in stats" :key="item.label">
        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
          <p class="mb-3 text-xs font-700 uppercase tracking-0.24em text-brand-700">{{ item.label }}</p>
          <NStatistic :value="item.value" />
          <p class="mt-3 text-sm text-[var(--app-text-muted)]">{{ item.caption }}</p>
        </NCard>
      </NGridItem>
    </NGrid>

    <div class="mt-6 grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Projects">
        <NDataTable
          :columns="projectColumns"
          :data="projects"
          :loading="loading"
          :pagination="{ pageSize: 5 }"
          :bordered="false"
        />
      </NCard>

      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Plugin registry">
        <NDataTable
          :columns="pluginColumns"
          :data="plugins"
          :loading="loading"
          :pagination="{ pageSize: 5 }"
          :bordered="false"
        />
      </NCard>
    </div>

    <NCard class="mt-6 rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Team footprint">
      <template v-if="teamMembers.length > 0">
        <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          <div
            v-for="user in teamMembers"
            :key="user.id"
            class="rounded-3xl border border-[var(--app-border)] bg-white/86 px-4 py-4"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <p class="text-sm font-700 text-ink-900">{{ user.name }}</p>
                <p class="mt-1 text-sm text-[var(--app-text-muted)]">{{ user.email }}</p>
              </div>
              <NTag
                round
                :bordered="false"
                :color="getRoleTagColor(user.role)"
              >
                {{ user.role }}
              </NTag>
            </div>
          </div>
        </div>
      </template>
      <NEmpty v-else description="No team members returned yet." />
    </NCard>
  </div>
</template>
