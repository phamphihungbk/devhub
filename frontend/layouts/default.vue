<script setup lang="ts">
const runtimeConfig = useRuntimeConfig()
</script>

<template>
  <div class="min-h-screen bg-[radial-gradient(circle_at_top,rgba(56,189,248,0.22),transparent_24%),linear-gradient(180deg,#f7fbff_0%,#eff6ff_35%,#f8fafc_100%)] text-slate-900 dark:bg-[radial-gradient(circle_at_top,rgba(14,165,233,0.18),transparent_20%),linear-gradient(180deg,#020617_0%,#0f172a_45%,#111827_100%)] dark:text-slate-100">
    <div class="mx-auto flex min-h-screen w-full max-w-7xl flex-col px-4 pb-10 pt-6 sm:px-6 lg:px-8">
      <header class="mb-8 flex flex-col gap-6 rounded-[2rem] border border-white/50 bg-white/80 px-6 py-6 shadow-[0_30px_80px_-40px_rgba(15,23,42,0.45)] backdrop-blur-xl dark:border-white/10 dark:bg-slate-900/70 sm:px-8">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div class="space-y-4">
            <UBadge color="primary" variant="soft" class="w-fit rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.28em]">
              Internal Developer Platform
            </UBadge>
            <div class="space-y-3">
              <p class="text-sm font-semibold uppercase tracking-[0.32em] text-sky-600 dark:text-sky-400">
                {{ runtimeConfig.public.appName }}
              </p>
              <h1 class="max-w-3xl text-4xl font-semibold tracking-tight text-slate-950 dark:text-white sm:text-5xl">
                Platform delivery, guardrails, and service visibility in one dashboard.
              </h1>
            </div>
            <p class="max-w-2xl text-sm leading-7 text-slate-600 dark:text-slate-300 sm:text-base">
              A Vue and Tailwind-first control plane with Nuxt UI primitives for projects, plugins, deployments, and developer workflows.
            </p>
          </div>

          <div class="grid gap-3 rounded-2xl border border-sky-100 bg-sky-50/80 p-4 text-sm text-slate-700 dark:border-sky-500/20 dark:bg-sky-500/10 dark:text-slate-200">
            <div class="flex items-center gap-2">
              <span class="h-2.5 w-2.5 rounded-full bg-emerald-500" />
              nginx is the shared gateway for the UI and API.
            </div>
            <div class="text-slate-500 dark:text-slate-400">
              Frontend: <span class="font-medium text-slate-800 dark:text-slate-100">devhub.local</span>
            </div>
            <div class="text-slate-500 dark:text-slate-400">
              API: <span class="font-medium text-slate-800 dark:text-slate-100">api.devhub.local</span>
            </div>
          </div>
        </div>
      </header>

      <main class="flex-1">
        <slot />
      </main>
    </div>
  </div>
</template>
