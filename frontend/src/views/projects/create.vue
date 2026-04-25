<script setup lang="ts">
import {
  NButton,
  NCard,
  NCheckbox,
  NCheckboxGroup,
  NForm,
  NFormItem,
  NInput,
  NSelect,
} from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useProjectCreateService } from '@/services/project'

const {
  environmentSelectOptions,
  form,
  openProjects,
  saving,
  scmProviderOptions,
  statusOptions,
  submitProject,
} = useProjectCreateService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Registry"
      title="Create project"
      description="Register a new service space in DevHub so the platform can track ownership, deployment targets, and future scaffolding workflows."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="openProjects">
          Back to projects
        </NButton>
        <NButton
          type="primary"
          :loading="saving"
          @click="submitProject"
        >
          Create project
        </NButton>
      </div>
    </PageHeader>

    <div class="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Project details">
        <NForm label-placement="top">
          <div class="grid gap-4 md:grid-cols-2">
            <NFormItem label="Project name" class="md:col-span-2">
              <NInput v-model:value="form.name" placeholder="payments-api" />
            </NFormItem>

            <NFormItem label="Status">
              <NSelect
                v-model:value="form.status"
                :options="statusOptions"
                placeholder="Select status"
              />
            </NFormItem>

            <NFormItem label="SCM provider">
              <NSelect
                v-model:value="form.scm_provider"
                :options="scmProviderOptions"
                placeholder="Select provider"
              />
            </NFormItem>

            <NFormItem label="Description" class="md:col-span-2">
              <NInput
                v-model:value="form.description"
                type="textarea"
                :autosize="{ minRows: 4, maxRows: 6 }"
                placeholder="Short summary of what this project owns and why it exists."
              />
            </NFormItem>
          </div>
        </NForm>
      </NCard>

      <div class="grid gap-6">
        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Deployment environments">
          <NFormItem label="Available environments">
            <NCheckboxGroup v-model:value="form.environments">
              <div class="grid gap-3">
                <NCheckbox
                  v-for="option in environmentSelectOptions"
                  :key="option.value"
                  :value="option.value"
                  :label="option.label"
                />
              </div>
            </NCheckboxGroup>
          </NFormItem>
          <p class="text-sm leading-6 text-[var(--app-text-muted)]">
            Choose the environments this service should support from day one. You can expand the lifecycle later as the platform grows.
          </p>
        </NCard>

        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="What gets registered">
          <div class="space-y-4 text-sm leading-6 text-[var(--app-text-muted)]">
            <p>
              DevHub will create a project record under your current team and the environments your operators can act on.
            </p>
            <p>
              This lays the groundwork for deployments, scaffold requests, release automation, and future control-plane workflows tied to this project.
            </p>
          </div>
        </NCard>
      </div>
    </div>
  </div>
</template>
