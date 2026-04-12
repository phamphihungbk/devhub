<script setup lang="ts">
import type { DashboardOverview, HighlightTone, StatusTone } from '~/types/dashboard'

const { data, pending, error, refresh } = await useFetch<DashboardOverview>('/api/internal/overview', {
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

const generatedAt = computed(() => {
  const value = data.value?.generatedAt
  if (!value) return 'Unavailable'

  return new Intl.DateTimeFormat('en', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
})

function toneToType(tone: HighlightTone) {
  switch (tone) {
    case 'accent':
      return 'info'
    case 'positive':
      return 'success'
    default:
      return 'default'
  }
}

function statusToType(tone: StatusTone) {
  switch (tone) {
    case 'accent':
      return 'info'
    case 'positive':
      return 'success'
    case 'warning':
      return 'warning'
    default:
      return 'default'
  }
}
</script>

<template>
  <div class="dashboard">
    <NCard class="naive-card" :bordered="false">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Overview</p>
          <h2>Platform snapshot</h2>
          <p class="section-copy">Service inventory, plugin registry, users, and recent delivery activity.</p>
        </div>
        <NButton type="primary" secondary :loading="pending" @click="refresh()">
          Refresh
        </NButton>
      </div>

      <NAlert
        v-if="error"
        type="error"
        title="Dashboard data request failed"
        class="empty-card"
      >
        Check that nginx, the frontend, and the backend are all running and that the API gateway is reachable.
      </NAlert>

      <div v-else class="stats-grid">
        <div v-for="highlight in highlights" :key="highlight.label" class="stat-card">
          <NTag size="small" :type="toneToType(highlight.tone)" round>
            {{ highlight.label }}
          </NTag>
          <p class="stat-card__value">{{ highlight.value }}</p>
          <p class="stat-card__description">{{ highlight.description }}</p>
        </div>
      </div>
    </NCard>

    <div class="dashboard-grid dashboard-grid--sidebar">
      <NCard class="naive-card" :bordered="false">
        <div class="section-heading">
          <div>
            <p class="section-heading__eyebrow">Projects</p>
            <h2>Service portfolio</h2>
            <p class="section-copy">Projects exposed by the backend through the gateway.</p>
          </div>
          <NTag round>{{ projects.length }} tracked</NTag>
        </div>

        <NEmpty v-if="projects.length === 0" description="No projects returned yet." class="empty-card" />

        <div v-else class="grid-cards grid-cards--projects">
          <article v-for="project in projects" :key="project.id" class="project-card">
            <div class="project-card__top">
              <div>
                <h3 class="project-card__title">{{ project.name }}</h3>
                <p class="project-card__body">
                  {{ project.description || 'No project description provided yet.' }}
                </p>
              </div>
              <NTag size="small" type="info" round>
                {{ project.environments.length }} envs
              </NTag>
            </div>

            <div class="tag-list">
              <NTag v-for="environment in project.environments" :key="environment" size="small" type="success" round>
                {{ environment }}
              </NTag>
            </div>
          </article>
        </div>
      </NCard>

      <div class="card-stack">
        <NCard class="naive-card" :bordered="false">
          <div class="section-heading">
            <div>
              <p class="section-heading__eyebrow">Gateway</p>
              <h2>Traffic summary</h2>
            </div>
          </div>

          <div class="summary-list">
            <div class="summary-item">
              <p class="summary-item__label">Frontend domain</p>
              <p class="summary-item__value">https://devhub.local</p>
            </div>
            <div class="summary-item">
              <p class="summary-item__label">API domain</p>
              <p class="summary-item__value">https://api.devhub.local</p>
            </div>
            <div class="summary-item">
              <p class="summary-item__label">Last snapshot</p>
              <p class="summary-item__value">{{ generatedAt }}</p>
            </div>
          </div>
        </NCard>

        <NCard class="naive-card" :bordered="false">
          <div class="section-heading">
            <div>
              <p class="section-heading__eyebrow">Plugins</p>
              <h2>Registry</h2>
            </div>
            <NTag round>{{ plugins.length }}</NTag>
          </div>

          <NEmpty v-if="plugins.length === 0" description="No plugins available." class="empty-card" />

          <div v-else class="card-stack">
            <article v-for="plugin in plugins" :key="plugin.id" class="resource-card">
              <div class="resource-card__top">
                <div>
                  <h3 class="resource-card__title">{{ plugin.name }}</h3>
                  <p class="resource-card__body">{{ plugin.scope }} • {{ plugin.version }}</p>
                </div>
                <NTag size="small" round>
                  {{ plugin.type }}
                </NTag>
              </div>
            </article>
          </div>
        </NCard>

        <NCard class="naive-card" :bordered="false">
          <div class="section-heading">
            <div>
              <p class="section-heading__eyebrow">Team</p>
              <h2>Known users</h2>
            </div>
            <NTag round>{{ users.length }}</NTag>
          </div>

          <NEmpty v-if="users.length === 0" description="No users available." class="empty-card" />

          <div v-else class="card-stack">
            <article v-for="user in users" :key="user.id" class="resource-card">
              <div class="resource-card__top">
                <div>
                  <h3 class="resource-card__title">{{ user.name }}</h3>
                  <p class="resource-card__body">{{ user.email }}</p>
                </div>
                <NTag size="small" type="info" round>
                  {{ user.role }}
                </NTag>
              </div>
            </article>
          </div>
        </NCard>
      </div>
    </div>

    <NCard class="naive-card" :bordered="false">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Deployments</p>
          <h2>Recent activity</h2>
          <p class="section-copy">Recent deployment records pulled from project-level backend endpoints.</p>
        </div>
        <NTag round>{{ deployments.length }} results</NTag>
      </div>

      <NEmpty v-if="deployments.length === 0" description="Deployments will appear here once services begin shipping." class="empty-card" />

      <div v-else class="grid-cards grid-cards--deployments">
        <article v-for="deployment in deployments" :key="deployment.id" class="deployment-card">
          <div class="deployment-card__top">
            <div>
              <h3 class="deployment-card__title">{{ deployment.service }}</h3>
              <p class="deployment-card__body">
                {{ deployment.environment }} • {{ deployment.version }}
              </p>
            </div>
            <NTag size="small" :type="statusToType(deployment.statusTone)" round>
              {{ deployment.status }}
            </NTag>
          </div>
          <p class="deployment-card__body">Project {{ deployment.projectId }}</p>
        </article>
      </div>
    </NCard>
  </div>
</template>
