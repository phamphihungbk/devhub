<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import {
  Logout,
  UserAvatar,
} from '@vicons/carbon'
import {
  NAvatar,
  NButton,
  NBreadcrumb,
  NBreadcrumbItem,
  NDropdown,
  NIcon,
  NLayout,
  NLayoutContent,
  NLayoutHeader,
  NLayoutSider,
  NMenu,
} from 'naive-ui'
import { computed, h, onMounted, ref } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'

import { renderRouteIcon } from '@/router'
import { useAuthStore } from '@/stores/modules/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const collapsed = ref(false)

onMounted(async () => {
  if (!authStore.profile && authStore.isAuthenticated) {
    try {
      await authStore.loadProfile()
    } catch {
      router.replace({ name: 'login' })
    }
  }
})

const menuOptions = computed<MenuOption[]>(() =>
  router.getRoutes()
    .filter((item) => item.name && item.meta.title)
    .filter((item) => !item.meta.guestOnly)
    .filter((item) => !item.meta.hideInMenu)
    .filter((item) => item.path !== '/')
    .map((item) => ({
      key: item.name as string,
      icon: renderRouteIcon(item.meta.icon as never),
      label: () => h(RouterLink, { to: { name: item.name as string } }, { default: () => item.meta.title as string }),
    })),
)

const breadcrumbItems = computed(() =>
  route.matched
    .filter((item) => item.meta?.title)
    .map((item) => ({
      key: item.path,
      title: item.meta.title as string,
    })),
)

const activeKey = computed(() => (route.meta.activeMenu as string) || (route.name as string) || 'dashboard')

const userMenu = computed<MenuOption[]>(() => [
  {
    key: 'role',
    label: authStore.profile?.role || 'Member',
    disabled: true,
    icon: () => h(NIcon, null, { default: () => h(UserAvatar) }),
  },
  {
    key: 'logout',
    label: 'Sign out',
    icon: () => h(NIcon, null, { default: () => h(Logout) }),
  },
])

async function handleUserAction(key: string) {
  if (key !== 'logout') return
  authStore.clearSession()
  await router.push({ name: 'login' })
}
</script>

<template>
  <NLayout has-sider class="min-h-screen bg-transparent text-ink-900">
    <NLayoutSider
      bordered
      collapse-mode="width"
      :collapsed="collapsed"
      :collapsed-width="84"
      :width="280"
      content-style="display:flex;flex-direction:column;padding:20px 16px 24px;"
      class="border-r border-[var(--app-border)] bg-[var(--app-sidebar)] backdrop-blur-xl"
    >
      <div class="mb-8 flex items-center gap-3 px-2">
        <div class="grid h-12 w-12 place-items-center rounded-2xl bg-[linear-gradient(145deg,#0f172a_0%,#1d4ed8_100%)] text-white shadow-[0_18px_38px_rgba(37,99,235,0.22)]">
          DH
        </div>
        <div v-if="!collapsed">
          <p class="text-xs font-600 uppercase tracking-0.24em text-ink-500">Platform control</p>
          <h1 class="text-lg font-700 text-ink-900">DevHub Console</h1>
        </div>
      </div>

      <div
        v-if="!collapsed"
        class="mb-5 rounded-3xl border border-[var(--app-border)] bg-[linear-gradient(180deg,rgba(255,255,255,0.96)_0%,rgba(244,247,251,0.9)_100%)] px-4 py-4 text-sm text-ink-700 shadow-[var(--app-shadow)]"
      >
        <p class="mb-1 text-xs font-700 uppercase tracking-0.22em text-brand-700">Control plane posture</p>
        <p class="leading-6">
          A route-first operator shell for service inventory, automation plugins, and access management across the platform stack.
        </p>
      </div>

      <NMenu
        :collapsed="collapsed"
        :collapsed-width="84"
        :collapsed-icon-size="22"
        :value="activeKey"
        :options="menuOptions"
      />

      <div class="mt-auto px-2 pt-6">
        <NButton class="w-full justify-start rounded-2xl" @click="collapsed = !collapsed">
          {{ collapsed ? 'Expand' : 'Collapse sidebar' }}
        </NButton>
      </div>
    </NLayoutSider>

    <NLayout>
      <NLayoutHeader bordered class="border-b border-[var(--app-border)] bg-[rgba(248,251,255,0.82)] px-8 py-5 backdrop-blur-xl">
        <div class="flex items-center justify-between gap-6">
          <div>
            <p class="text-xs font-700 uppercase tracking-0.24em text-brand-700">Internal developer platform</p>
            <NBreadcrumb class="mt-2">
              <NBreadcrumbItem v-for="item in breadcrumbItems" :key="item.key">
                {{ item.title }}
              </NBreadcrumbItem>
            </NBreadcrumb>
          </div>

          <div class="flex items-center gap-3">
            <div class="hidden text-right sm:block">
              <p class="text-sm font-700 text-ink-900">{{ authStore.profile?.name || 'Operator' }}</p>
              <p class="text-xs text-ink-500">{{ authStore.profile?.email || 'Signed session' }}</p>
            </div>

            <NDropdown :options="userMenu" @select="handleUserAction">
              <button class="flex items-center gap-3 rounded-2xl border border-[var(--app-border)] bg-white/88 px-3 py-2 text-left shadow-sm transition hover:border-brand-300 hover:shadow-md">
                <NAvatar round color="#2563eb">
                  {{ (authStore.profile?.name || 'D').slice(0, 1).toUpperCase() }}
                </NAvatar>
                <span class="hidden text-sm font-600 text-ink-700 sm:inline">Session</span>
              </button>
            </NDropdown>
          </div>
        </div>
      </NLayoutHeader>

      <NLayoutContent class="min-h-[calc(100vh-88px)] bg-[radial-gradient(circle_at_top_right,_rgba(37,99,235,0.1),_transparent_26%),linear-gradient(180deg,_rgba(248,251,255,0.92)_0%,_rgba(237,242,247,0.98)_100%)] px-6 py-6 md:px-8">
        <RouterView />
      </NLayoutContent>
    </NLayout>
  </NLayout>
</template>
