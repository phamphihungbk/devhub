<script setup lang="ts">
import { NButton, NCard, NForm, NFormItem, NInput, useMessage } from 'naive-ui'
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { ApiError } from '@/services/request'
import { useAuthStore } from '@/stores/modules/auth'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const message = useMessage()

const form = reactive({
  email: '',
  password: '',
})
const pending = ref(false)

async function submit() {
  if (!form.email || !form.password) {
    message.warning('Email and password are required.')
    return
  }

  pending.value = true

  try {
    await authStore.login(form)
    message.success('Welcome back.')
    await router.replace((route.query.redirect as string) || '/dashboard')
  } catch (error) {
    const fallback = 'Unable to sign in with the provided credentials.'
    message.error(error instanceof ApiError ? error.message : fallback)
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <main class="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(37,99,235,0.12),_transparent_24%),linear-gradient(145deg,_#f6f9fc_0%,_#edf2f7_50%,_#e2e8f0_100%)] px-6 py-8">
    <div class="mx-auto grid min-h-[calc(100vh-64px)] max-w-6xl gap-8 lg:grid-cols-[1.1fr_0.9fr]">
      <section class="flex items-center">
        <div class="max-w-xl">
          <p class="text-xs font-800 uppercase tracking-0.28em text-brand-700">DevHub Control Plane</p>
          <h1 class="mt-5 text-5xl font-900 leading-tight tracking-tight text-ink-900">
            A steadier operator console for your internal platform.
          </h1>
          <p class="mt-6 text-lg leading-8 text-[var(--app-text-muted)]">
            DevHub gives platform teams one place to inspect services, manage automation, and operate the delivery surface with less noise.
          </p>
          <div class="mt-8 grid gap-4 sm:grid-cols-2">
            <div class="rounded-3xl border border-white/60 bg-white/70 p-5 shadow-sm backdrop-blur">
              <p class="text-sm font-700 text-ink-900">Project visibility</p>
              <p class="mt-2 text-sm leading-6 text-[var(--app-text-muted)]">
                Browse service inventory, environments, and the operational context behind each workspace.
              </p>
            </div>
            <div class="rounded-3xl border border-white/60 bg-white/70 p-5 shadow-sm backdrop-blur">
              <p class="text-sm font-700 text-ink-900">Platform operations</p>
              <p class="mt-2 text-sm leading-6 text-[var(--app-text-muted)]">
                Review plugin automation and access controls from a cleaner, more focused control surface.
              </p>
            </div>
          </div>
        </div>
      </section>

      <section class="flex items-center justify-center">
        <NCard
          class="w-full max-w-lg rounded-[28px] border border-[var(--app-border)] shadow-[0_30px_80px_rgba(15,23,42,0.12)]"
          content-style="padding: 32px;"
        >
          <div class="mb-8">
            <p class="text-xs font-800 uppercase tracking-0.24em text-brand-700">Sign in</p>
            <h2 class="mt-3 text-3xl font-800 text-ink-900">Enter the control plane</h2>
            <p class="mt-2 text-sm leading-6 text-[var(--app-text-muted)]">
              Use a DevHub account to access service inventory, plugin automation, and the operator workflows behind the platform.
            </p>
          </div>

          <NForm label-placement="top" @submit.prevent="submit">
            <NFormItem label="Email">
              <NInput v-model:value="form.email" placeholder="operator@devhub.local" />
            </NFormItem>

            <NFormItem label="Password">
              <NInput v-model:value="form.password" type="password" show-password-on="click" placeholder="••••••••" />
            </NFormItem>

            <NButton
              type="primary"
              attr-type="submit"
              block
              size="large"
              class="h-12 rounded-2xl text-sm font-700 shadow-[0_18px_32px_rgba(37,99,235,0.22)]"
            >
              Continue
            </NButton>
          </NForm>
        </NCard>
      </section>
    </div>
  </main>
</template>
