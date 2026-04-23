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
  useMessage,
} from 'naive-ui'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import PageHeader from '@/components/page-header.vue'
import { createProject } from '@/services/api'
import { ApiError } from '@/services/request'
import { useAuthStore } from '@/stores/modules/auth'
import { environmentOptions } from '@/theme/environment'
import type { ProjectPayload } from '@/services/api'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()
const saving = ref(false)

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Active', value: 'active' },
  { label: 'Archived', value: 'archived' },
  { label: 'Deprecated', value: 'deprecated' },
]

const scmProviderOptions = [
  { label: 'Gitea', value: 'gitea' },
  { label: 'GitHub', value: 'github' },
  { label: 'GitLab', value: 'gitlab' },
  { label: 'Bitbucket', value: 'bitbucket' },
]

const form = reactive<ProjectPayload>({
  name: '',
  description: '',
  environments: ['dev'],
  status: 'draft',
  team_id: '',
  scm_provider: 'gitea',
})

function validateForm() {
  if (!form.name.trim()) return 'Project name is required.'
  if (form.environments.length === 0) return 'Select at least one environment.'
  if (!authStore.profile?.team_id) return 'Your team information is not available.'
  if (!form.scm_provider.trim()) return 'SCM provider is required.'
  return null
}

async function submit() {
  const validationError = validateForm()

  if (validationError) {
    message.warning(validationError)
    return
  }

  saving.value = true

  try {
    await createProject({
      ...form,
      team_id: authStore.profile.team_id,
      name: form.name.trim(),
      description: form.description?.trim() || undefined,
      scm_provider: form.scm_provider.trim(),
    })

    message.success('Project created successfully.')
    await router.push({ name: 'projects' })
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to create project.')
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Registry"
      title="Create project"
      description="Register a new service space in DevHub so the platform can track ownership, deployment targets, and future scaffolding workflows."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="router.push({ name: 'projects' })">
          Back to projects
        </NButton>
        <NButton
          type="primary"
          :loading="saving"
          @click="submit"
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
                  v-for="option in environmentOptions"
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
