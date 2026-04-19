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
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, ref } from 'vue'

import PageHeader from '@/components/page-header.vue'
import { fetchPlugins, fetchProjects, fetchUsers } from '@/services/api'
import { ApiError } from '@/services/request'
import { getEnvironmentTagColor } from '@/theme/environment'
import { getPluginTypeTagColor } from '@/theme/plugin'
import type { PluginRecord, Project, UserRecord } from '@/services/api'

const message = useMessage()
const loading = ref(false)
const projects = ref<Project[]>([])
const plugins = ref<PluginRecord[]>([])
const users = ref<UserRecord[]>([])

const stats = computed(() => [
  {
    label: 'Projects',
    value: projects.value.length,
    caption: 'Tracked service domains',
  },
  {
    label: 'Plugins',
    value: plugins.value.length,
    caption: 'Installed platform automations',
  },
  {
    label: 'Users',
    value: users.value.length,
    caption: 'Operators with console access',
  },
])

const projectColumns = [
  {
    title: 'Project',
    key: 'name',
  },
  {
    title: 'Environments',
    key: 'environments',
    render: (row: Project) =>
      h(
        'div',
        { class: 'flex flex-wrap gap-2' },
        row.environments.map((value) =>
          h(
            NTag,
            {
              bordered: false,
              color: getEnvironmentTagColor(value),
            },
            { default: () => value },
          ),
        ),
      ),
  },
  {
    title: 'Description',
    key: 'description',
    render: (row: Project) => row.description || 'No description yet.',
  },
]

const pluginColumns = [
  { title: 'Plugin', key: 'name' },
  {
    title: 'Type',
    key: 'type',
    render: (row: PluginRecord) =>
      h(
        NTag,
        {
          bordered: false,
          color: getPluginTypeTagColor(row.type),
        },
        { default: () => row.type },
      ),
  },
  { title: 'Runtime', key: 'runtime' },
  { title: 'Scope', key: 'scope' },
]

async function load() {
  loading.value = true
  try {
    const [projectData, pluginData, userData] = await Promise.all([
      fetchProjects(),
      fetchPlugins(),
      fetchUsers(),
    ])
    projects.value = projectData
    plugins.value = pluginData
    users.value = userData
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load dashboard data.')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Overview"
      title="Platform dashboard"
      description="A higher-signal entry point for the control plane: service inventory, automation registry, and operator visibility in one calmer admin workspace."
    >
      <NButton type="primary" @click="load">
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
      <template v-if="users.length > 0">
        <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          <div
            v-for="user in users"
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
                :color="{ color: '#e0f2fe', textColor: '#0369a1' }"
              >
                {{ user.role }}
              </NTag>
            </div>
          </div>
        </div>
      </template>
      <NEmpty v-else description="No users returned yet." />
    </NCard>
  </div>
</template>
