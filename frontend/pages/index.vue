<script setup lang="ts">
import type { DashboardOverview } from '~/types/dashboard'

const { data, pending, error, refresh } = await useFetch<DashboardOverview>('/api/overview', {
  default: () => ({
    generatedAt: new Date().toISOString(),
    highlights: [],
    projects: [],
    plugins: [],
    users: [],
    deployments: [],
  }),
})

const highlights = computed(() => data.value?.highlights ?? [])
const projects = computed(() => data.value?.projects ?? [])
const plugins = computed(() => data.value?.plugins ?? [])
const users = computed(() => data.value?.users ?? [])
const deployments = computed(() => data.value?.deployments ?? [])
</script>

<template>
  <div class="dashboard">
    <section class="panel panel--highlights">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Overview</p>
          <h2>Platform snapshot</h2>
        </div>
        <button class="button" type="button" @click="refresh()">
          Refresh
        </button>
      </div>

      <div v-if="pending" class="empty-state">
        Loading platform data...
      </div>
      <div v-else-if="error" class="empty-state empty-state--error">
        <p>Frontend reached the Nuxt server, but the backend data request failed.</p>
        <p class="empty-state__hint">Check that the API is running on the configured backend URL.</p>
      </div>
      <div v-else class="metric-grid">
        <MetricCard
          v-for="highlight in highlights"
          :key="highlight.label"
          :label="highlight.label"
          :value="highlight.value"
          :description="highlight.description"
          :tone="highlight.tone"
        />
      </div>
    </section>

    <section class="panel panel--split">
      <div>
        <div class="section-heading">
          <div>
            <p class="section-heading__eyebrow">Projects</p>
            <h2>Service portfolio</h2>
          </div>
          <span class="chip">{{ projects.length }} tracked</span>
        </div>
        <div v-if="projects.length === 0" class="empty-state">
          No projects returned yet.
        </div>
        <div v-else class="project-list">
          <article v-for="project in projects" :key="project.id" class="project-card">
            <div class="project-card__header">
              <h3>{{ project.name }}</h3>
              <span class="chip chip--muted">{{ project.environments.length }} envs</span>
            </div>
            <p>{{ project.description || 'No project description provided yet.' }}</p>
            <div class="tag-row">
              <span v-for="environment in project.environments" :key="environment" class="tag">
                {{ environment }}
              </span>
            </div>
          </article>
        </div>
      </div>

      <div class="stack">
        <ResourceTable
          title="Plugins"
          subtitle="Extensibility registry"
          :rows="plugins"
          empty-message="No plugins available."
        >
          <template #default="{ row }">
            <div class="table-row__primary">
              <span>{{ row.name }}</span>
              <StatusBadge :label="row.type" tone="neutral" />
            </div>
            <div class="table-row__secondary">
              <span>{{ row.version }}</span>
              <span>{{ row.scope }}</span>
            </div>
          </template>
        </ResourceTable>

        <ResourceTable
          title="Team"
          subtitle="Known platform users"
          :rows="users"
          empty-message="No users available."
        >
          <template #default="{ row }">
            <div class="table-row__primary">
              <span>{{ row.name }}</span>
              <StatusBadge :label="row.role" tone="accent" />
            </div>
            <div class="table-row__secondary">
              <span>{{ row.email }}</span>
            </div>
          </template>
        </ResourceTable>
      </div>
    </section>

    <section class="panel">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Deployments</p>
          <h2>Recent activity</h2>
        </div>
        <span class="section-heading__meta">{{ deployments.length }} results</span>
      </div>
      <div v-if="deployments.length === 0" class="empty-state">
        Deployments will appear here once projects begin shipping.
      </div>
      <div v-else class="deployment-list">
        <article v-for="deployment in deployments" :key="deployment.id" class="deployment-card">
          <div class="deployment-card__header">
            <div>
              <h3>{{ deployment.service }}</h3>
              <p>{{ deployment.environment }} • {{ deployment.version }}</p>
            </div>
            <StatusBadge :label="deployment.status" :tone="deployment.statusTone" />
          </div>
          <p class="deployment-card__meta">Project {{ deployment.projectId }}</p>
        </article>
      </div>
    </section>
  </div>
</template>
