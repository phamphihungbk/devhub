<script setup lang="ts">
import type { DashboardOverview, HighlightTone, StatusTone } from '~/types/dashboard'

const { data, pending, error, refresh } = await useFetch<DashboardOverview>('/internal/overview', {
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
  if (!value) {
    return 'Unavailable'
  }

  return new Intl.DateTimeFormat('en', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
})

function toneToColor(tone: HighlightTone) {
  switch (tone) {
    case 'accent':
      return 'primary'
    case 'positive':
      return 'success'
    default:
      return 'neutral'
  }
}

function statusToneToColor(tone: StatusTone) {
  switch (tone) {
    case 'accent':
      return 'primary'
    case 'positive':
      return 'success'
    case 'warning':
      return 'warning'
    default:
      return 'neutral'
  }
}
</script>

<template>
  <div class="space-y-6">
    <section class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_22rem]">
      <UCard class="overflow-hidden border-0 bg-slate-950 text-white shadow-2xl shadow-slate-950/10">
        <div class="space-y-6">
          <div class="flex flex-wrap items-start justify-between gap-4">
            <div class="space-y-2">
              <p class="text-sm font-medium uppercase tracking-[0.24em] text-sky-300">
                Dashboard overview
              </p>
              <h2 class="text-2xl font-semibold tracking-tight sm:text-3xl">
                Platform health, projects, plugins, and delivery flow.
              </h2>
              <p class="max-w-2xl text-sm leading-7 text-slate-300">
                This dashboard is fed by the backend through nginx, giving the frontend one stable gateway for internal traffic.
              </p>
            </div>

            <UButton color="primary" variant="solid" icon="i-lucide-refresh-cw" :loading="pending" @click="refresh()">
              Refresh data
            </UButton>
          </div>

          <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
            <div
              v-for="highlight in highlights"
              :key="highlight.label"
              class="rounded-2xl border border-white/10 bg-white/5 p-4 backdrop-blur-sm"
            >
              <UBadge :color="toneToColor(highlight.tone)" variant="soft" class="mb-4 rounded-full px-2.5 py-1">
                {{ highlight.label }}
              </UBadge>
              <div class="text-3xl font-semibold tracking-tight">
                {{ highlight.value }}
              </div>
              <p class="mt-2 text-sm leading-6 text-slate-300">
                {{ highlight.description }}
              </p>
            </div>
          </div>
        </div>
      </UCard>

      <UCard class="border-slate-200/70 bg-white/85 backdrop-blur dark:border-slate-800 dark:bg-slate-900/75">
        <div class="space-y-5">
          <div>
            <p class="text-sm font-semibold uppercase tracking-[0.24em] text-sky-600 dark:text-sky-400">
              Gateway
            </p>
            <h2 class="mt-2 text-xl font-semibold text-slate-950 dark:text-white">
              Traffic summary
            </h2>
          </div>

          <div class="space-y-4 text-sm text-slate-600 dark:text-slate-300">
            <div class="rounded-2xl bg-slate-50 p-4 dark:bg-slate-800/80">
              <div class="font-medium text-slate-900 dark:text-slate-100">Frontend domain</div>
              <div class="mt-1">https://devhub.local</div>
            </div>
            <div class="rounded-2xl bg-slate-50 p-4 dark:bg-slate-800/80">
              <div class="font-medium text-slate-900 dark:text-slate-100">API domain</div>
              <div class="mt-1">https://api.devhub.local</div>
            </div>
            <div class="rounded-2xl bg-sky-50 p-4 text-sky-950 dark:bg-sky-500/10 dark:text-sky-100">
              <div class="font-medium">Last snapshot</div>
              <div class="mt-1">{{ generatedAt }}</div>
            </div>
          </div>
        </div>
      </UCard>
    </section>

    <section v-if="error" class="rounded-3xl border border-rose-200 bg-rose-50/90 p-6 text-rose-900 dark:border-rose-500/20 dark:bg-rose-500/10 dark:text-rose-100">
      <p class="text-base font-semibold">Dashboard data request failed.</p>
      <p class="mt-2 text-sm text-rose-700 dark:text-rose-200">
        Check that nginx, the frontend, and the backend are running and that the API gateway is reachable.
      </p>
    </section>

    <section class="grid gap-6 xl:grid-cols-[minmax(0,1.35fr)_minmax(0,1fr)]">
      <UCard class="border-slate-200/70 bg-white/85 backdrop-blur dark:border-slate-800 dark:bg-slate-900/75">
        <div class="space-y-5">
          <div class="flex items-center justify-between gap-4">
            <div>
              <p class="text-sm font-semibold uppercase tracking-[0.24em] text-sky-600 dark:text-sky-400">
                Projects
              </p>
              <h2 class="mt-2 text-xl font-semibold text-slate-950 dark:text-white">Service portfolio</h2>
            </div>
            <UBadge color="neutral" variant="soft" class="rounded-full px-3 py-1">
              {{ projects.length }} tracked
            </UBadge>
          </div>

          <div v-if="projects.length === 0" class="rounded-2xl border border-dashed border-slate-200 p-8 text-center text-sm text-slate-500 dark:border-slate-700 dark:text-slate-400">
            No projects returned by the backend yet.
          </div>

          <div v-else class="grid gap-4 md:grid-cols-2">
            <article
              v-for="project in projects"
              :key="project.id"
              class="rounded-2xl border border-slate-200 bg-slate-50/70 p-5 transition hover:-translate-y-0.5 hover:border-sky-200 hover:shadow-lg hover:shadow-sky-100/50 dark:border-slate-800 dark:bg-slate-800/60 dark:hover:border-sky-500/30 dark:hover:shadow-sky-950/20"
            >
              <div class="flex items-start justify-between gap-4">
                <div>
                  <h3 class="text-base font-semibold text-slate-950 dark:text-white">{{ project.name }}</h3>
                  <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-slate-300">
                    {{ project.description || 'No project description provided yet.' }}
                  </p>
                </div>
                <UBadge color="secondary" variant="soft" class="rounded-full px-2.5 py-1">
                  {{ project.environments.length }} envs
                </UBadge>
              </div>

              <div class="mt-4 flex flex-wrap gap-2">
                <UBadge
                  v-for="environment in project.environments"
                  :key="environment"
                  color="primary"
                  variant="subtle"
                  class="rounded-full px-2.5 py-1"
                >
                  {{ environment }}
                </UBadge>
              </div>
            </article>
          </div>
        </div>
      </UCard>

      <div class="space-y-6">
        <UCard class="border-slate-200/70 bg-white/85 backdrop-blur dark:border-slate-800 dark:bg-slate-900/75">
          <div class="space-y-4">
            <div class="flex items-center justify-between gap-4">
              <div>
                <p class="text-sm font-semibold uppercase tracking-[0.24em] text-sky-600 dark:text-sky-400">
                  Plugins
                </p>
                <h2 class="mt-2 text-xl font-semibold text-slate-950 dark:text-white">Registry</h2>
              </div>
              <UBadge color="neutral" variant="soft" class="rounded-full px-3 py-1">{{ plugins.length }}</UBadge>
            </div>

            <div v-if="plugins.length === 0" class="rounded-2xl border border-dashed border-slate-200 p-6 text-sm text-slate-500 dark:border-slate-700 dark:text-slate-400">
              No plugins available.
            </div>

            <div v-else class="space-y-3">
              <article
                v-for="plugin in plugins"
                :key="plugin.id"
                class="flex items-start justify-between gap-4 rounded-2xl border border-slate-200 bg-slate-50/70 p-4 dark:border-slate-800 dark:bg-slate-800/60"
              >
                <div class="min-w-0">
                  <div class="truncate font-medium text-slate-950 dark:text-white">{{ plugin.name }}</div>
                  <div class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                    {{ plugin.scope }} • {{ plugin.version }}
                  </div>
                </div>
                <UBadge color="neutral" variant="soft" class="rounded-full px-2.5 py-1">{{ plugin.type }}</UBadge>
              </article>
            </div>
          </div>
        </UCard>

        <UCard class="border-slate-200/70 bg-white/85 backdrop-blur dark:border-slate-800 dark:bg-slate-900/75">
          <div class="space-y-4">
            <div class="flex items-center justify-between gap-4">
              <div>
                <p class="text-sm font-semibold uppercase tracking-[0.24em] text-sky-600 dark:text-sky-400">
                  Team
                </p>
                <h2 class="mt-2 text-xl font-semibold text-slate-950 dark:text-white">Known users</h2>
              </div>
              <UBadge color="neutral" variant="soft" class="rounded-full px-3 py-1">{{ users.length }}</UBadge>
            </div>

            <div v-if="users.length === 0" class="rounded-2xl border border-dashed border-slate-200 p-6 text-sm text-slate-500 dark:border-slate-700 dark:text-slate-400">
              No users available.
            </div>

            <div v-else class="space-y-3">
              <article
                v-for="user in users"
                :key="user.id"
                class="flex items-start justify-between gap-4 rounded-2xl border border-slate-200 bg-slate-50/70 p-4 dark:border-slate-800 dark:bg-slate-800/60"
              >
                <div class="min-w-0">
                  <div class="truncate font-medium text-slate-950 dark:text-white">{{ user.name }}</div>
                  <div class="mt-1 truncate text-sm text-slate-500 dark:text-slate-400">{{ user.email }}</div>
                </div>
                <UBadge color="primary" variant="soft" class="rounded-full px-2.5 py-1">{{ user.role }}</UBadge>
              </article>
            </div>
          </div>
        </UCard>
      </div>
    </section>

    <section>
      <UCard class="border-slate-200/70 bg-white/85 backdrop-blur dark:border-slate-800 dark:bg-slate-900/75">
        <div class="space-y-5">
          <div class="flex items-center justify-between gap-4">
            <div>
              <p class="text-sm font-semibold uppercase tracking-[0.24em] text-sky-600 dark:text-sky-400">
                Deployments
              </p>
              <h2 class="mt-2 text-xl font-semibold text-slate-950 dark:text-white">Recent activity</h2>
            </div>
            <UBadge color="neutral" variant="soft" class="rounded-full px-3 py-1">{{ deployments.length }} results</UBadge>
          </div>

          <div v-if="deployments.length === 0" class="rounded-2xl border border-dashed border-slate-200 p-8 text-center text-sm text-slate-500 dark:border-slate-700 dark:text-slate-400">
            Deployments will appear here once services begin shipping.
          </div>

          <div v-else class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <article
              v-for="deployment in deployments"
              :key="deployment.id"
              class="rounded-2xl border border-slate-200 bg-slate-50/70 p-5 dark:border-slate-800 dark:bg-slate-800/60"
            >
              <div class="flex items-start justify-between gap-4">
                <div>
                  <h3 class="text-base font-semibold text-slate-950 dark:text-white">{{ deployment.service }}</h3>
                  <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                    {{ deployment.environment }} • {{ deployment.version }}
                  </p>
                </div>
                <UBadge :color="statusToneToColor(deployment.statusTone)" variant="soft" class="rounded-full px-2.5 py-1">
                  {{ deployment.status }}
                </UBadge>
              </div>
              <p class="mt-4 text-sm text-slate-600 dark:text-slate-300">
                Project {{ deployment.projectId }}
              </p>
            </article>
          </div>
        </div>
      </UCard>
    </section>
  </div>
</template>
