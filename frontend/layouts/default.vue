<script setup lang="ts">
const runtimeConfig = useRuntimeConfig()
</script>

<template>
  <div class="shell">
    <div class="shell__backdrop shell__backdrop--top" />
    <div class="shell__backdrop shell__backdrop--bottom" />
    <div class="shell__grid" />
    <main class="shell__content">
      <header class="hero">
        <div class="hero__eyebrow">Internal Developer Platform</div>
        <div class="hero__heading">
          <p class="hero__kicker">{{ runtimeConfig.public.appName }}</p>
          <h1>Ship services without losing sight of the system.</h1>
        </div>
        <p class="hero__copy">
          A Nuxt-powered control surface for projects, plugins, and delivery flow. The backend stays in charge;
          this frontend gives it a clearer voice.
        </p>
      </header>
      <div class="hero__meta">
        <span class="hero__pill">Frontend: devhub.local</span>
        <span class="hero__pill">API: api.devhub.local</span>
        <span class="hero__pill">Gateway: nginx</span>
      </div>
      <slot />
    </main>
  </div>
</template>
