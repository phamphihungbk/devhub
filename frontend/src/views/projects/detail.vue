<script setup lang="ts">
import {
  NButton,
  NCard,
  NDataTable,
  NStatistic,
  NTag,
} from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useProjectDetailService } from '@/services/project'

const {
  deploymentColumns,
  deployments,
  failedDeployments,
  failedReleases,
  getEnvironmentTagColor,
  loading,
  openProjects,
  openService,
  ownerContact,
  ownerTeamName,
  project,
  releaseColumns,
  releases,
  serviceColumns,
  services,
  successfulDeployments,
  successfulReleases,
} = useProjectDetailService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Registry"
      :title="project ? project.name : 'Project details'"
      description="Review the services under this project and the most recent lifecycle activity before deciding what to release or deploy next."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="openProjects">
          Back to projects
        </NButton>
      </div>
    </PageHeader>

    <div class="grid gap-4 md:grid-cols-5">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Services" :value="services.length" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Successful releases" :value="successfulReleases" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Failed releases" :value="failedReleases" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Successful deployments" :value="successfulDeployments" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Failed deployments" :value="failedDeployments" />
      </NCard>
    </div>

    <div class="mt-6 grid gap-6 xl:grid-cols-[0.95fr_1.05fr]">
      <div class="grid gap-6">
        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Project posture">
          <div class="grid gap-4 text-sm leading-6 text-[var(--app-text-muted)] md:grid-cols-2">
            <div class="md:col-span-2">
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Environments
              </p>
              <div class="mt-3 flex flex-wrap gap-2">
                <NTag
                  v-for="environment in project?.environments || []"
                  :key="environment"
                  :bordered="false"
                  :color="getEnvironmentTagColor(environment)"
                >
                  {{ environment }}
                </NTag>
              </div>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Status
              </p>
              <p class="mt-1 text-base font-semibold text-[var(--app-text)]">
                {{ project?.status || 'Unknown' }}
              </p>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Owner team
              </p>
              <p class="mt-1 text-base font-semibold text-[var(--app-text)]">
                {{ ownerTeamName }}
              </p>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                SCM provider
              </p>
              <p class="mt-1 text-base font-semibold text-[var(--app-text)]">
                {{ project?.scm_provider || 'Not set' }}
              </p>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Owner contact
              </p>
              <p class="mt-1 text-base font-semibold text-[var(--app-text)]">
                {{ ownerContact }}
              </p>
            </div>
          </div>
          <p class="mt-4 text-sm leading-6 text-[var(--app-text-muted)]">
            {{ project?.description || 'No description provided yet.' }}
          </p>
        </NCard>

        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Services">
          <NDataTable
            :columns="serviceColumns"
            :data="services"
            :loading="loading"
            :pagination="{ pageSize: 6 }"
            :bordered="false"
            :row-props="row => ({
              class: 'cursor-pointer',
              onClick: () => openService(row),
            })"
          />
        </NCard>
      </div>

      <div class="grid gap-6">
        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Recent releases">
          <NDataTable
            :columns="releaseColumns"
            :data="releases"
            :loading="loading"
            :pagination="{ pageSize: 6 }"
            :bordered="false"
          />
        </NCard>

        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Recent deployments">
          <NDataTable
            :columns="deploymentColumns"
            :data="deployments"
            :loading="loading"
            :pagination="{ pageSize: 6 }"
            :bordered="false"
          />
        </NCard>
      </div>
    </div>

  </div>
</template>
